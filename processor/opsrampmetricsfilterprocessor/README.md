# Alert Metrics Extractor Processor

The Alert Metrics Extractor Processor is a filtering processor that dynamically extracts raw metric names from PromQL alert expressions stored in Kubernetes ConfigMaps and filters incoming metrics to only pass those that are referenced in alert definitions.

## Description

This processor reads alert definitions from a Kubernetes ConfigMap, parses PromQL expressions to extract raw metric names, maintains a thread-safe global map of distinct metrics, and filters incoming metrics to only allow those present in the map. It watches for ConfigMap changes and updates the filtering rules in real-time.

## Configuration

```yaml
processors:
  alertmetricsextractor:
    # Name of the ConfigMap containing alert definitions
    alert_definitions_configmap_name: "opsramp-agent-config"
    
    # Key in the ConfigMap containing the alert definitions YAML
    alert_definitions_configmap_key: "alert-definitions.yaml"
    
    # Kubernetes namespace where the ConfigMap is located
    namespace: "default"
```

### Configuration Parameters

| Parameter | Type | Default | Required | Description |
|-----------|------|---------|----------|-------------|
| `alert_definitions_configmap_name` | string | `opsramp-agent-config` | Yes | Name of the ConfigMap containing alert definitions |
| `alert_definitions_configmap_key` | string | `alert-definitions.yaml` | Yes | Key in the ConfigMap containing alert definitions YAML |
| `namespace` | string | `default` | Yes | Kubernetes namespace where the ConfigMap is located |

## How It Works

### Processing Flow

1. **Initialization**
   - Creates Kubernetes client to access ConfigMaps
   - Reads initial alert definitions from the specified ConfigMap
   - Starts background goroutine to watch for ConfigMap changes

2. **Metric Extraction**
   - Parses YAML alert definitions from ConfigMap
   - For each alert rule, extracts the PromQL expression (`expr` field)
   - Uses Prometheus PromQL parser to walk the AST and identify metric names
   - Handles both `VectorSelector` (instant metrics) and `MatrixSelector` (range metrics)

3. **Global Map Building**
   - Creates a distinct set of all extracted metric names
   - Thread-safely updates the global filtering map
   - Logs the total number of distinct metrics found

4. **Metric Filtering**
   - For each incoming metric batch in the pipeline:
     - Checks if metric name exists in the global map
     - **Allows through**: Metrics present in alert definitions
     - **Filters out**: Metrics not referenced in any alert
   - Passes filtered metrics to the next consumer/exporter

5. **Real-time Updates**
   - Continuously watches the ConfigMap for changes
   - When ConfigMap is modified, re-extracts metrics and updates the global map
   - All subsequent filtering uses the updated rules

### Alert Definitions Format

The processor expects alert definitions in the following YAML structure:

```yaml
alertDefinitions:
  - resourceType: "k8s_cluster"
    rules:
      - name: "api_server_availability"
        expr: "(sum(increase(apiserver_request_total{verb!=\"WATCH\"}[5m])) / sum(increase(apiserver_request_total{verb!=\"WATCH\"}[5m]))) * 100"
        # ... other alert fields
      
  - resourceType: "k8s_node"
    rules:
      - name: "node_memory_usage"
        expr: "(k8s_node_memory_working_set / k8s_node_memory_available) * 100"
        # ... other alert fields
```

### PromQL Metric Extraction

The processor can extract metrics from various PromQL constructs:

- **Simple metrics**: `up`, `cpu_usage_percent`
- **Metrics with labels**: `apiserver_request_total{verb="GET"}`
- **Function calls**: `rate(http_requests_total[5m])`, `increase(errors_total[1h])`
- **Range queries**: `node_cpu_seconds_total[5m]`
- **Complex expressions**: `(1 - (node_memory_available / node_memory_total)) * 100`
- **Aggregations**: `sum(rate(http_requests_total[5m]))`

#### Example Extraction

From this PromQL expression:
```promql
(sum(increase(apiserver_request_total{verb!="WATCH"}[5m])) / sum(increase(apiserver_request_total{verb!="WATCH"}[5m]))) * 100
```

The processor extracts: `["apiserver_request_total"]`

From this expression:
```promql
rate(http_requests_total[5m]) + rate(http_errors_total[5m])
```

The processor extracts: `["http_requests_total", "http_errors_total"]`

## Usage Examples

### Basic Pipeline Configuration

```yaml
receivers:
  prometheus:
    config:
      scrape_configs:
        - job_name: 'kubernetes-metrics'
          kubernetes_sd_configs:
            - role: node

processors:
  alertmetricsextractor:
    alert_definitions_configmap_name: "opsramp-agent-config"
    alert_definitions_configmap_key: "alert-definitions.yaml"
    namespace: "monitoring"
  
  batch:
    timeout: 1s

exporters:
  prometheusremotewrite:
    endpoint: "https://prometheus.example.com/api/v1/write"

service:
  pipelines:
    metrics/filtered:
      receivers: [prometheus]
      processors: [alertmetricsextractor, batch]
      exporters: [prometheusremotewrite]
```

### ConfigMap Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: opsramp-agent-config
  namespace: monitoring
data:
  alert-definitions.yaml: |
    alertDefinitions:
      - resourceType: k8s_cluster
        rules:
          - name: k8s_apiserver_requests_error_rate
            expr: (sum(increase(apiserver_request_total{verb!="WATCH",code=~"2.."}[5m])) / sum(increase(apiserver_request_total{verb!="WATCH"}[5m]))) * 100
          - name: k8s_cluster_nodes_health
            expr: (sum(k8s_node_condition_ready) / count(k8s_node_condition_ready)) * 100
      - resourceType: k8s_pod
        rules:
          - name: k8s_pod_cpu_usage
            expr: k8s_pod_cpu_limit_utilization_ratio * 100
```

## Thread Safety

The processor ensures thread-safe operation through:

- **Read-Write Mutex**: Uses `sync.RWMutex` for the global metrics map
- **Concurrent Reads**: Multiple goroutines can read the map simultaneously during metric filtering
- **Exclusive Writes**: ConfigMap updates acquire exclusive write access to update the map
- **Atomic Operations**: Map updates are atomic to prevent race conditions

## Error Handling

| Error Type | Behavior | Log Level |
|------------|----------|-----------|
| ConfigMap not found | Processor continues with empty map (filters all metrics) | Error |
| Invalid YAML | Processor continues with previous map | Error |
| PromQL parse error | Skips that expression, continues with others | Warning |
| Kubernetes API error | Retries with exponential backoff | Error |
| ConfigMap deleted | Clears metrics map (filters all metrics) | Warning |

## Performance Considerations

- **Memory Usage**: Only stores metric names (strings), not metric data
- **CPU Impact**: PromQL parsing done only when ConfigMap changes, not per metric
- **Filtering Performance**: O(1) map lookup per metric name
- **Network**: Minimal - only watches ConfigMap changes via Kubernetes API
- **Scalability**: Handles thousands of metrics efficiently

## Troubleshooting

### Common Issues

1. **No metrics passing through**
   - Check if ConfigMap exists and is accessible
   - Verify alert definitions contain valid PromQL expressions
   - Check processor logs for parsing errors

2. **ConfigMap changes not detected**
   - Ensure proper RBAC permissions for ConfigMap access
   - Check if ConfigMap is in the correct namespace
   - Verify Kubernetes watch connection

3. **Some metrics still filtered**
   - Check if metric names match exactly (case-sensitive)
   - Verify PromQL expressions reference the expected metrics
   - Enable debug logging to see extracted metric names

### Debug Logging

Enable debug logging to see detailed information:

```yaml
service:
  telemetry:
    logs:
      level: debug
```

This will show:
- Extracted metric names from each alert expression
- ConfigMap watch events
- Filtering decisions for each metric

## RBAC Requirements

The processor requires the following Kubernetes permissions:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otelcol-alertmetricsextractor
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "list", "watch"]
```