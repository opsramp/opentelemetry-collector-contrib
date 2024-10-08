package redis

import (
	"context"
	"encoding/json"
	"time"

	lru "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/lru"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Client struct {
	Host     string
	Port     string
	Password string
	GoClient *goredis.Client
	lruCache *lru.Cache
	Enabled  bool
	logger   *zap.Logger
}

type OpsrampRedisConfig struct {
	RedisHost         string        `mapstructure:"redisHost"`
	RedisPort         string        `mapstructure:"redisPort"`
	RedisPass         string        `mapstructure:"redisPass"`
	ClusterName       string        `mapstructure:"clusterName"`
	ClusterUid        string        `mapstructure:"clusterUid"`
	NodeName          string        `mapstructure:"nodeName"`
	LruCacheSize      int           `mapstructure:"lruCacheSize"`
	LruExpirationTime time.Duration `mapstructure:"lruExpirationTime"`
}

func NewClient(logger *zap.Logger, lruCache *lru.Cache, rHost, rPort, rPass string) *Client {
	client := Client{
		Host:     rHost,
		Port:     rPort,
		Password: rPass,
		Enabled:  true,
		logger:   logger,
		lruCache: lruCache,
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

func (c *Client) GetUuidValueInString(ctx context.Context, key string) string {
	logger := c.logger

	// Try to init the cache if it is firt time

	cache := c.lruCache

	value, ok := cache.Get(key)
	if !ok {
		logger.Debug("Failed to fetch the key from lru cache ", zap.Any("key", key))
		if c.Enabled {
			val, err := c.GoClient.Get(ctx, key).Result()
			if err == goredis.Nil {
				logger.Debug("key does not exist ", zap.Any("key", key))
			} else if err != nil {
				logger.Error("Failed to fetch the key from redis ", zap.Error(err))
			} else {
				logger.Debug("Got value from redis ", zap.Any("key", key), zap.Any("value", value))
				type RedisData struct {
					ResourceUuid string `json:"resourceUuid,omitempty"`
					ResourceHash uint64 `json:"resourceHash,omitempty"`
				}
				var redisData RedisData
				err = json.Unmarshal([]byte(val), &redisData)
				if err != nil {
					logger.Error("Could not unmarshal data:", zap.Error(err))
					return ""
				}
				value = redisData.ResourceUuid
			}
		}

		if value != "" {
			cache.Add(key, value)
		}
	} else {
		logger.Debug("Got value from lru cache ", zap.Any("key", key), zap.Any("value", value))
	}
	return value
}
