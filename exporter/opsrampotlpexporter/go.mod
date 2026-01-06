module github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampotlpexporter

go 1.24.0

require (
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/common v0.104.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.104.0
	github.com/opsramp/go-proxy-dialer v0.0.0-20240313152735-64bb1ce65640
	github.com/stretchr/testify v1.11.1
	go.opentelemetry.io/collector/component/componenttest v0.141.0
	go.opentelemetry.io/collector/config/configauth v1.47.0
	go.opentelemetry.io/collector/config/configcompression v1.47.0
	go.opentelemetry.io/collector/config/configgrpc v0.141.0
	go.opentelemetry.io/collector/config/configopaque v1.47.0
	go.opentelemetry.io/collector/config/configretry v1.47.0
	go.opentelemetry.io/collector/config/configtls v1.47.0
	go.opentelemetry.io/collector/confmap v1.47.0
	go.opentelemetry.io/collector/consumer v1.47.0
	go.opentelemetry.io/collector/consumer/consumererror v0.141.0
	go.opentelemetry.io/collector/exporter v1.47.0
	go.opentelemetry.io/collector/exporter/exporterhelper v0.141.0
	go.opentelemetry.io/collector/exporter/exportertest v0.141.0
	go.opentelemetry.io/collector/otelcol/otelcoltest v0.141.0
	go.opentelemetry.io/collector/pdata v1.47.0
	go.uber.org/goleak v1.3.0
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8
	google.golang.org/grpc v1.77.0
	google.golang.org/protobuf v1.36.10
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/ebitengine/purego v0.9.0 // indirect
	github.com/foxboron/go-tpm-keyfiles v0.0.0-20250903184740-5d135037bd4d // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/google/go-tpm v0.9.7 // indirect
	github.com/grafana/regexp v0.0.0-20240518133315-a468a5bfb3bc // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.2 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/knadh/koanf/v2 v2.3.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/otlptranslator v0.0.2 // indirect
	github.com/shirou/gopsutil/v4 v4.25.10 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/collector/client v1.47.0 // indirect
	go.opentelemetry.io/collector/component/componentstatus v0.141.0 // indirect
	go.opentelemetry.io/collector/config/configmiddleware v1.47.0 // indirect
	go.opentelemetry.io/collector/config/confignet v1.47.0 // indirect
	go.opentelemetry.io/collector/config/configoptional v1.47.0 // indirect
	go.opentelemetry.io/collector/config/configtelemetry v0.141.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/envprovider v1.47.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/fileprovider v1.47.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/httpprovider v1.47.0 // indirect
	go.opentelemetry.io/collector/confmap/provider/yamlprovider v1.47.0 // indirect
	go.opentelemetry.io/collector/confmap/xconfmap v0.141.0 // indirect
	go.opentelemetry.io/collector/connector v0.141.0 // indirect
	go.opentelemetry.io/collector/connector/connectortest v0.141.0 // indirect
	go.opentelemetry.io/collector/connector/xconnector v0.141.0 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.141.0 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.141.0 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.141.0 // indirect
	go.opentelemetry.io/collector/extension v1.47.0 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v1.47.0 // indirect
	go.opentelemetry.io/collector/extension/extensioncapabilities v0.141.0 // indirect
	go.opentelemetry.io/collector/extension/extensionmiddleware v0.141.0 // indirect
	go.opentelemetry.io/collector/extension/extensiontest v0.141.0 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.141.0 // indirect
	go.opentelemetry.io/collector/featuregate v1.47.0 // indirect
	go.opentelemetry.io/collector/internal/fanoutconsumer v0.141.0 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.141.0 // indirect
	go.opentelemetry.io/collector/otelcol v0.141.0 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.141.0 // indirect
	go.opentelemetry.io/collector/pdata/testdata v0.141.0 // indirect
	go.opentelemetry.io/collector/pdata/xpdata v0.141.0 // indirect
	go.opentelemetry.io/collector/pipeline v1.47.0 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.141.0 // indirect
	go.opentelemetry.io/collector/processor v1.47.0 // indirect
	go.opentelemetry.io/collector/processor/processortest v0.141.0 // indirect
	go.opentelemetry.io/collector/processor/xprocessor v0.141.0 // indirect
	go.opentelemetry.io/collector/receiver v1.47.0 // indirect
	go.opentelemetry.io/collector/receiver/receivertest v0.141.0 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.141.0 // indirect
	go.opentelemetry.io/collector/service v0.141.0 // indirect
	go.opentelemetry.io/collector/service/hostcapabilities v0.141.0 // indirect
	go.opentelemetry.io/contrib/bridges/otelzap v0.13.0 // indirect
	go.opentelemetry.io/contrib/otelconf v0.18.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutlog v0.14.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.38.0 // indirect
	go.opentelemetry.io/otel/log v0.14.0 // indirect
	go.opentelemetry.io/otel/sdk/log v0.14.0 // indirect
	go.opentelemetry.io/proto/otlp v1.7.1 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.4 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	gonum.org/v1/gonum v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.1 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20220913051719-115f729f3c8c // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.3-0.20250322232337-35a7c28c31ee // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_golang v1.23.2 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.1 // indirect
	github.com/prometheus/procfs v0.17.0 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/collector/component v1.47.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.63.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.38.0 // indirect
	go.opentelemetry.io/otel v1.38.0 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.60.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk v1.38.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/net v0.47.0
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	go.opentelemetry.io/collector => ../../../opentelemetry-collector
	go.opentelemetry.io/collector/client => ../../../opentelemetry-collector/client
	go.opentelemetry.io/collector/component => ../../../opentelemetry-collector/component
	go.opentelemetry.io/collector/component/componentstatus => ../../../opentelemetry-collector/component/componentstatus
	go.opentelemetry.io/collector/config/configauth => ../../../opentelemetry-collector/config/configauth
	go.opentelemetry.io/collector/config/configcompression => ../../../opentelemetry-collector/config/configcompression
	go.opentelemetry.io/collector/config/configgrpc => ../../../opentelemetry-collector/config/configgrpc
	go.opentelemetry.io/collector/config/confignet => ../../../opentelemetry-collector/config/confignet
	go.opentelemetry.io/collector/config/configopaque => ../../../opentelemetry-collector/config/configopaque
	go.opentelemetry.io/collector/config/configretry => ../../../opentelemetry-collector/config/configretry
	go.opentelemetry.io/collector/config/configtelemetry => ../../../opentelemetry-collector/config/configtelemetry
	go.opentelemetry.io/collector/config/configtls => ../../../opentelemetry-collector/config/configtls
	go.opentelemetry.io/collector/confmap => ../../../opentelemetry-collector/confmap
	go.opentelemetry.io/collector/confmap/provider/envprovider => ../../../opentelemetry-collector/confmap/provider/envprovider
	go.opentelemetry.io/collector/confmap/provider/fileprovider => ../../../opentelemetry-collector/confmap/provider/fileprovider
	go.opentelemetry.io/collector/confmap/provider/httpprovider => ../../../opentelemetry-collector/confmap/provider/httpprovider
	go.opentelemetry.io/collector/confmap/provider/yamlprovider => ../../../opentelemetry-collector/confmap/provider/yamlprovider
	go.opentelemetry.io/collector/connector => ../../../opentelemetry-collector/connector
	go.opentelemetry.io/collector/connector/connectortest => ../../../opentelemetry-collector/connector/connectortest
	go.opentelemetry.io/collector/consumer => ../../../opentelemetry-collector/consumer
	go.opentelemetry.io/collector/consumer/consumererror => ../../../opentelemetry-collector/consumer/consumererror
	go.opentelemetry.io/collector/consumer/consumertest => ../../../opentelemetry-collector/consumer/consumertest
	go.opentelemetry.io/collector/exporter => ../../../opentelemetry-collector/exporter
	go.opentelemetry.io/collector/exporter/exportertest => ../../../opentelemetry-collector/exporter/exportertest
	go.opentelemetry.io/collector/extension => ../../../opentelemetry-collector/extension
	go.opentelemetry.io/collector/extension/auth => ../../../opentelemetry-collector/extension/extensionauth
	go.opentelemetry.io/collector/extension/extensioncapabilities => ../../../opentelemetry-collector/extension/extensioncapabilities
	go.opentelemetry.io/collector/featuregate => ../../../opentelemetry-collector/featuregate
	go.opentelemetry.io/collector/internal/fanoutconsumer => ../../../opentelemetry-collector/internal/fanoutconsumer
	go.opentelemetry.io/collector/otelcol => ../../../opentelemetry-collector/otelcol
	go.opentelemetry.io/collector/otelcol/otelcoltest => ../../../opentelemetry-collector/otelcol/otelcoltest
	go.opentelemetry.io/collector/pdata => ../../../opentelemetry-collector/pdata
	go.opentelemetry.io/collector/pdata/pprofile => ../../../opentelemetry-collector/pdata/pprofile
	go.opentelemetry.io/collector/pdata/testdata => ../../../opentelemetry-collector/pdata/testdata
	go.opentelemetry.io/collector/pipeline => ../../../opentelemetry-collector/pipeline
	go.opentelemetry.io/collector/processor => ../../../opentelemetry-collector/processor
	go.opentelemetry.io/collector/processor/processortest => ../../../opentelemetry-collector/processor/processortest
	go.opentelemetry.io/collector/receiver => ../../../opentelemetry-collector/receiver
	go.opentelemetry.io/collector/receiver/receivertest => ../../../opentelemetry-collector/receiver/receivertest
	go.opentelemetry.io/collector/service => ../../../opentelemetry-collector/service
)
