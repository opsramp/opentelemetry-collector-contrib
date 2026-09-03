package redis

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/cache"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Client struct {
	Host                       string
	Port                       string
	Password                   string
	GoClient                   *goredis.Client
	CacheObject                *cache.Cache
	PrimaryCacheEvictionTime   time.Duration
	SecondaryCacheEvictionTime time.Duration
	Enabled                    bool
	logger                     *zap.Logger
}

type OpsrampRedisConfig struct {
	RedisHost                  string        `mapstructure:"redisHost"`
	RedisPort                  string        `mapstructure:"redisPort"`
	RedisPass                  string        `mapstructure:"redisPass"`
	ClusterName                string        `mapstructure:"clusterName"`
	ClusterUid                 string        `mapstructure:"clusterUid"`
	NodeName                   string        `mapstructure:"nodeName"`
	PrimaryCacheSize           int           `mapstructure:"primaryCacheSize"`
	SecondaryCacheSize         int           `mapstructure:"secondaryCacheSize"`
	PrimaryCacheEvictionTime   time.Duration `mapstructure:"primaryCacheEvictionTime"`
	SecondaryCacheEvictionTime time.Duration `mapstructure:"secondaryCacheEvictionTime"`
	EnableGpuNicRouting        bool          `mapstructure:"enableGpuNicRouting"`

	// EnrichAttributesFromRedis backfills the k8s.* labels stored alongside the
	// resource uuid. Intended for passthrough mode, where no informers run.
	EnrichAttributesFromRedis bool `mapstructure:"enrichAttributesFromRedis"`
}

func NewClient(logger *zap.Logger, cache *cache.Cache, rHost, rPort, rPass string, primaryCacheEvictionTime, secondaryCacheEvictionTime time.Duration) *Client {
	client := Client{
		Host:                       rHost,
		Port:                       rPort,
		Password:                   rPass,
		Enabled:                    true,
		logger:                     logger,
		CacheObject:                cache,
		PrimaryCacheEvictionTime:   primaryCacheEvictionTime,
		SecondaryCacheEvictionTime: secondaryCacheEvictionTime,
	}

	if client.Host == "" {
		logger.Info("Redis Host is empty, hence no lookup for moid/resourceuuid cache")
		client.Enabled = false
		return &client
	}

	if client.Port == "" {
		client.Port = "6379"
	}

	client.Init()

	return &client
}

func (c *Client) Init() error {
	c.GoClient = goredis.NewClient(&goredis.Options{
		Addr:            c.Host + ":" + c.Port,
		Password:        c.Password,
		MaxRetries:      -1,
		MinRetryBackoff: 55 * time.Millisecond,
		MaxRetryBackoff: 2 * time.Second,
	})

	if err := c.TestConnection(context.Background()); err != nil {
		return err
	}

	return nil
}

func (c *Client) TestConnection(ctx context.Context) error {
	logger := c.logger

	var err error
	err = nil

	for i := 0; i < 15; i++ {
		_, err = c.GoClient.Ping(ctx).Result()
		if err != nil {
			logger.Info("Could not connect/ping to Redis", zap.Any("error", err.Error()))
		} else {
			logger.Info("Connected to Redis")
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		logger.Error("Could not connect/ping to Redis", zap.Any("error", err.Error()))
	}

	return err
}

// RedisData mirrors the payload the OpsRamp agent writes under each MoID key.
// Field names and json tags must stay in sync with the agent's cache.RedisData.
type RedisData struct {
	ResourceUuid    string `json:"resourceUuid,omitempty"`
	ResourceHash    uint64 `json:"resourceHash,omitempty"`
	NodeName        string `json:"k8s.node.name,omitempty"`
	NamespaceName   string `json:"k8s.namespace.name,omitempty"`
	DeploymentName  string `json:"k8s.deployment.name,omitempty"`
	DaemonSetName   string `json:"k8s.daemonset.name,omitempty"`
	StatefulSetName string `json:"k8s.statefulset.name,omitempty"`
	ReplicaSetName  string `json:"k8s.replicaset.name,omitempty"`
	PodName         string `json:"k8s.pod.name,omitempty"`
	PodUid          string `json:"k8s.pod.uid,omitempty"`
	PodIp           string `json:"k8s.pod.ip,omitempty"`
}

// GetResourceUuid is nil-safe so callers can test a lookup result inline.
func (d *RedisData) GetResourceUuid() string {
	if d == nil {
		return ""
	}
	return d.ResourceUuid
}

// GetRedisData is the single lookup path for a MoID key: at most one cache read
// and one Redis GET per call. Both cache tiers hold the raw JSON so a cache hit
// can serve the resource uuid and the labels without a second round trip.
func (c *Client) GetRedisData(ctx context.Context, key string) *RedisData {
	logger := c.logger

	if c.CacheObject != nil {
		if val, err := c.CacheObject.GetFromPrimary(key); err == nil {
			logger.Debug("Got value from PrimaryCache", zap.Any("key", key), zap.Any("value", val))
			return c.decode(key, val)
		}
		if val, err := c.CacheObject.GetFromSecondary(key); err == nil {
			logger.Debug("Got value from SecondaryCache", zap.Any("key", key), zap.Any("value", val))
			return c.decode(key, val)
		}
	}

	if !c.Enabled {
		return nil
	}

	val, err := c.GoClient.Get(ctx, key).Result()
	if err != nil {
		if err == goredis.Nil {
			logger.Debug("Key does not exist in Redis", zap.Any("key", key))
		} else {
			logger.Error("Failed to fetch the key from Redis", zap.Error(err))
		}
		c.negativeCache(key)
		return nil
	}

	logger.Debug("Got value from Redis", zap.Any("key", key), zap.Any("value", val))

	redisData := c.decode(key, val)
	if redisData == nil {
		c.negativeCache(key)
		return nil
	}

	if c.CacheObject != nil {
		if redisData.ResourceUuid == "" {
			c.CacheObject.AddToSecondaryWithTTL(key, val, c.SecondaryCacheEvictionTime)
		} else {
			c.CacheObject.AddToPrimaryWithTTL(key, val, c.PrimaryCacheEvictionTime)
		}
	}
	return redisData
}

// decode turns a cached/Redis payload into RedisData. An empty payload is the
// negative-cache marker and is reported as a miss. A non-JSON payload is a
// legacy cache entry holding the bare resource uuid.
func (c *Client) decode(key, val string) *RedisData {
	if val == "" {
		return nil
	}
	if !strings.HasPrefix(strings.TrimSpace(val), "{") {
		return &RedisData{ResourceUuid: val}
	}
	var redisData RedisData
	if err := json.Unmarshal([]byte(val), &redisData); err != nil {
		c.logger.Error("Could not unmarshal data", zap.Any("key", key), zap.Error(err))
		return nil
	}
	return &redisData
}

func (c *Client) negativeCache(key string) {
	if c.CacheObject != nil {
		c.CacheObject.AddToSecondaryWithTTL(key, "", c.SecondaryCacheEvictionTime)
	}
}

func (c *Client) GetUuidValueInString(ctx context.Context, key string) string {
	if redisData := c.GetRedisData(ctx, key); redisData != nil {
		return redisData.ResourceUuid
	}
	return ""
}

type RedisDataWithDeployment struct {
	ResourceUuid   string `json:"resourceUuid,omitempty"`
	ResourceHash   uint64 `json:"resourceHash,omitempty"`
	DeploymentName string `json:"k8s.deployment.name,omitempty"`
}

func (c *Client) GetRedisDataWithDeployment(ctx context.Context, key string) *RedisDataWithDeployment {
	redisData := c.GetRedisData(ctx, key)
	if redisData == nil {
		return nil
	}
	return &RedisDataWithDeployment{
		ResourceUuid:   redisData.ResourceUuid,
		ResourceHash:   redisData.ResourceHash,
		DeploymentName: redisData.DeploymentName,
	}
}
