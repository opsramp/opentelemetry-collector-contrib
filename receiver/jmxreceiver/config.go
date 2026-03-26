// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package jmxreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/jmxreceiver"

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/scraper/scraperhelper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// jmxGathererMainClass the class containing the main function for the JMX Metric Gatherer JAR
var jmxGathererMainClass = "io.opentelemetry.contrib.jmxmetrics.JmxMetrics"

// jmxScraperMainClass the class containing the main function for the JMX Scraper JAR
var jmxScraperMainClass = "io.opentelemetry.contrib.jmxscraper.JmxScraper"

type Config struct {
	Applications       map[string]ApplicationConfig `mapstructure:"applications"`
	OTLPExporterConfig otlpExporterConfig           `mapstructure:"otlp"`
}

type ApplicationConfig struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`

	// The path for the JMX Metric Gatherer or JMX Scraper JAR (/opt/opentelemetry-java-contrib-jmx-metrics.jar by default).
	// Supported by: jmx-scraper and jmx-metric-gatherer
	JARPath string `mapstructure:"jar_path"`
	// The Service URL or host:port for the target coerced to one of form: service:jmx:rmi:///jndi/rmi://<host>:<port>/jmxrmi.
	// Supported by: jmx-scraper and jmx-metric-gatherer
	Endpoint string `mapstructure:"endpoint"`
	// Comma-separated list of systems to monitor
	// Supported by: jmx-scraper and jmx-metric-gatherer
	TargetSystem string `mapstructure:"target_system"`
	// The target source of metric definitions to use for the target system.
	// Supported values are: auto, instrumentation and legacy.
	// Supported by: jmx-scraper
	TargetSource string `mapstructure:"target_source"`
	// Comma-separated list of paths to custom YAML metrics definition,
	// mandatory when TargetSystem is not set.
	// Supported by: jmx-scraper
	JmxConfigs string `mapstructure:"jmx_configs"`
	// The duration in between groovy script invocations and metric exports (10 seconds by default).
	// Will be converted to milliseconds.
	// Supported by: jmx-scraper and jmx-metric-gatherer
	CollectionInterval time.Duration `mapstructure:"collection_interval"`
	// The OTLP exporter settings
	// Supported by: jmx-scraper and jmx-metric-gatherer
	OTLPExporterConfig otlpExporterConfig `mapstructure:"otlp"`
	// The JMX username
	// Supported by: jmx-scraper and jmx-metric-gatherer
	Username string `mapstructure:"username"`
	// The JMX password
	// Supported by: jmx-scraper and jmx-metric-gatherer
	Password configopaque.String `mapstructure:"password"`
	// The keystore path for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	KeystorePath string `mapstructure:"keystore_path"`
	// The keystore password for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	KeystorePassword configopaque.String `mapstructure:"keystore_password"`
	// The keystore type for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	KeystoreType string `mapstructure:"keystore_type"`
	// The truststore path for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	TruststorePath string `mapstructure:"truststore_path"`
	// The truststore password for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	TruststorePassword configopaque.String `mapstructure:"truststore_password"`
	// The truststore type for SSL
	// Supported by: jmx-scraper and jmx-metric-gatherer
	TruststoreType string `mapstructure:"truststore_type"`
	// The JMX remote profile.  Should be one of:
	// `"SASL/PLAIN"`, `"SASL/DIGEST-MD5"`, `"SASL/CRAM-MD5"`, `"TLS SASL/PLAIN"`, `"TLS SASL/DIGEST-MD5"`, or
	// `"TLS SASL/CRAM-MD5"`, though no enforcement is applied.
	// Supported by: jmx-scraper and jmx-metric-gatherer
	RemoteProfile string `mapstructure:"remote_profile"`
	// The SASL/DIGEST-MD5 realm
	// Supported by: jmx-scraper and jmx-metric-gatherer
	Realm string `mapstructure:"realm"`
	// Array of additional JARs to be added to the class path when launching the JMX Metric Gatherer JAR
	// Supported by: jmx-scraper and jmx-metric-gatherer
	AdditionalJars []string `mapstructure:"additional_jars"`
	// Map of resource attributes used by the Java SDK Autoconfigure to set resource attributes
	// Supported by: jmx-scraper and jmx-metric-gatherer
	ResourceAttributes map[string]string `mapstructure:"resource_attributes"`
	// Log level used by the JMX metric gatherer. Should be one of:
	// `"trace"`, `"debug"`, `"info"`, `"warn"`, `"error"`, `"off"`
	// Supported by: jmx-metric-gatherer
	LogLevel string `mapstructure:"log_level"`
}

// We don't embed the existing OTLP Exporter config as most fields are unsupported
type otlpExporterConfig struct {
	// The OTLP Receiver endpoint to send metrics to ("0.0.0.0:<random open port>" by default).
	Endpoint string `mapstructure:"endpoint"`
	// The OTLP exporter timeout (5 seconds by default).  Will be converted to milliseconds.
	TimeoutSettings exporterhelper.TimeoutConfig `mapstructure:",squash"`
	// The headers to include in OTLP metric submission requests.
	Headers map[string]string `mapstructure:"headers"`
}

func (oec otlpExporterConfig) headersToString() string {
	// sort for reliable testing
	headers := make([]string, 0, len(oec.Headers))
	for k := range oec.Headers {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	headerString := ""
	for _, k := range headers {
		v := oec.Headers[k]
		headerString += fmt.Sprintf("%s=%v,", k, v)
	}
	// remove trailing comma
	headerString = headerString[0 : len(headerString)-1]
	return headerString
}

func (c *ApplicationConfig) parseProperties(logger *zap.Logger) []string {
	// slf4j.simpleLogger only available in JMX Metrics Gatherer jar
	//if err := c.validateJar(jmxMetricsGathererVersions, c.JARPath); err == nil {
	parsed := make([]string, 0, 1)

	logLevel := "info"
	if c.LogLevel != "" {
		logLevel = strings.ToLower(c.LogLevel)
	} else if logger != nil {
		logLevel = getZapLoggerLevelEquivalent(logger)
	}

	parsed = append(parsed, "-Dorg.slf4j.simpleLogger.defaultLogLevel="+logLevel)

	// Sorted for testing and reproducibility
	sort.Strings(parsed)
	return parsed
	//}
	//return nil
}

var logLevelTranslator = map[zapcore.Level]string{
	zap.DebugLevel:  "debug",
	zap.InfoLevel:   "info",
	zap.WarnLevel:   "warn",
	zap.ErrorLevel:  "error",
	zap.DPanicLevel: "error",
	zap.PanicLevel:  "error",
	zap.FatalLevel:  "error",
}

var zapLevels = []zapcore.Level{
	zap.DebugLevel,
	zap.InfoLevel,
	zap.WarnLevel,
	zap.ErrorLevel,
	zap.DPanicLevel,
	zap.PanicLevel,
	zap.FatalLevel,
}

func getZapLoggerLevelEquivalent(logger *zap.Logger) string {
	var loggerLevel *zapcore.Level
	for i, level := range zapLevels {
		if testLevel(logger, level) {
			loggerLevel = &zapLevels[i]
			break
		}
	}

	// Couldn't get log level from logger default logger level to info
	if loggerLevel == nil {
		return "info"
	}

	return logLevelTranslator[*loggerLevel]
}

func testLevel(logger *zap.Logger, level zapcore.Level) bool {
	return logger.Check(level, "_") != nil
}

// parseClasspath creates a classpath string with the JMX Gatherer JAR at the beginning
func (c *ApplicationConfig) parseClasspath() string {
	var classPathElems []string

	// Add JMX JAR to classpath
	classPathElems = append(classPathElems, c.JARPath)

	// Add additional JARs if any
	classPathElems = append(classPathElems, c.AdditionalJars...)

	// Join them
	return strings.Join(classPathElems, ":")
}

func isSupportedJAR(supportedJarDetails map[string]supportedJar, jar string) bool {
	hash, err := hashFile(jar)
	if err != nil {
		return false
	}
	_, ok := supportedJarDetails[hash]
	return ok
}

/*func (c *ApplicationConfig) jarMainClass() string {
	if isSupportedJAR(jmxMetricsGathererVersions, c.JARPath) {
		return jmxGathererMainClass
	} else if isSupportedJAR(jmxScraperVersions, c.JARPath) {
		return jmxScraperMainClass
	}
	return ""
}

func (c *ApplicationConfig) jarJMXSamplingConfig() (string, string) {
	if isSupportedJAR(jmxMetricsGathererVersions, c.JARPath) {
		return "otel.jmx.interval.milliseconds", strconv.FormatInt(c.CollectionInterval.Milliseconds(), 10)
	} else if isSupportedJAR(jmxScraperVersions, c.JARPath) {
		return "otel.metric.export.interval", c.CollectionInterval.String()
	}
	return "", ""
}*/

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (c *Config) validateJar(supportedJarDetails map[string]supportedJar, jar string) error {
	hash, err := hashFile(jar)
	if err != nil {
		return fmt.Errorf("error hashing file: %w", err)
	}

	jarDetails, ok := supportedJarDetails[hash]
	if !ok {
		return errors.New("jar hash does not match known versions")
	}
	if jarDetails.addedValidation != nil {
		if err = jarDetails.addedValidation(c, jarDetails); err != nil {
			return fmt.Errorf("jar failed validation: %w", err)
		}
	}

	return nil
}

var (
	validLogLevels = map[string]struct{}{"trace": {}, "debug": {}, "info": {}, "warn": {}, "error": {}, "off": {}}
	// feature parity between jmx-gatherer and jmx-scraper
	validTargetSystems = map[string]struct{}{
		"activemq": {}, "cassandra": {}, "hbase": {}, "hadoop": {},
		"jetty": {}, "jvm": {}, "kafka": {}, "kafka-consumer": {}, "kafka-producer": {}, "solr": {}, "tomcat": {}, "wildfly": {},
	}
)
var AdditionalTargetSystems = "n/a"

// Separated into two functions for tests
func init() {
	initAdditionalTargetSystems()
}

func initAdditionalTargetSystems() {
	if AdditionalTargetSystems != "n/a" {
		for t := range strings.SplitSeq(AdditionalTargetSystems, ",") {
			validTargetSystems[t] = struct{}{}
		}
	}
}

func (c *Config) Validate() error {
	if len(c.Applications) == 0 {
		return fmt.Errorf("at least one application must be configured")
	}

	for appName, appConfig := range c.Applications {
		if appConfig.JARPath == "" {
			return fmt.Errorf("missing required field `jar_path` for application %s", appName)
		}
		if appConfig.Endpoint == "" {
			return fmt.Errorf("missing required field `endpoint` for application %s", appName)
		}
		if appConfig.TargetSystem == "" {
			return fmt.Errorf("missing required field `target_system` for application %s", appName)
		}
		if appConfig.CollectionInterval < 0 {
			return fmt.Errorf("`collection_interval` must be positive for application %s: %vms", appName, appConfig.CollectionInterval.Milliseconds())
		}
		if len(appConfig.LogLevel) > 0 {
			if _, ok := validLogLevels[strings.ToLower(appConfig.LogLevel)]; !ok {
				return fmt.Errorf("`log_level` must be one of %s for application %s", listKeys(validLogLevels), appName)
			}
		}
	}

	return nil
}

func listKeys(presenceMap map[string]struct{}) string {
	list := make([]string, 0, len(presenceMap))
	for k := range presenceMap {
		list = append(list, fmt.Sprintf("'%s'", k))
	}
	sort.Strings(list)
	return strings.Join(list, ", ")
}
