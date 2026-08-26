# opsrampmetricsfilter — Metric Categories & Shared Definitions Loader

Reference for the `opsrampmetricsfilter` processor after the metric-categories
change: what was added, how the component is structured, and every supported
deployment pattern.

- **Component type:** `opsrampmetricsfilter`
- **Stability:** alpha
- **Signal:** metrics only
- **Source:** `processor/opsrampmetricsfilterprocessor/`

---

## 1. What the processor does

It is an **allow-list metrics filter derived from OpsRamp alert definitions**.
Rather than maintaining a hand-written list of metrics to keep, it parses the
PromQL `expr` of every alert rule, extracts the referenced metric names, and
drops every metric not in that set.

```
alert-definitions.yaml ──▶ PromQL parse ──▶ allow-list ──▶ filter pipeline metrics
```

---

## 2. Summary of changes

### 2.1 New capability — `metric_categories`

Alert definitions are now **categorized by `resourceType`**, and a processor
instance can be restricted to a subset of categories. This lets one set of alert
definitions drive multiple pipelines with different exporters.

Previously `resourceType` was parsed into `AlertDefinition.ResourceType` and then
discarded during flattening. It is now retained and drives categorization.

### 2.2 New architecture — shared definitions loader

All processor instances resolving to the same alert-definitions source share a
single refcounted loader. The ConfigMap watch, the YAML parse, and the PromQL
parse happen **once** regardless of how many instances exist.

Before (per instance):

```
instance ──▶ own k8s client ──▶ own watch ──▶ own parse ──▶ own allow-list
```

After (per source):

```
                                       ┌─▶ instance A (podMetric)
source ──▶ loader ──▶ parse ──▶ publish┤
          (1 watch)                    └─▶ instance B (clusterMetric)
```

Benefits:

| | Before (N instances) | After (N instances) |
|---|---|---|
| ConfigMap watches / API connections | N | 1 |
| YAML + PromQL parses per reload | N | 1 |
| Kubernetes clients | N | 1 |
| Reload skew between instances | possible | none |

### 2.3 File layout

| File | Status | Contents |
|---|---|---|
| `category.go` | **new** | `categorySet` bitmask, `resourceType` → category mapping, config value parsing |
| `loader.go` | **new** | Shared refcounted loader, source keying, YAML/PromQL parsing, ConfigMap watch |
| `filewatcher.go` | **new** | `FileWatcher` (fsnotify + polling fallback), moved out of `processor.go` |
| `category_test.go` | **new** | Category parsing, categorization, projection, multi-subscriber tests |
| `processor.go` | rewritten | Reduced to filtering + loader subscription (−586 lines) |
| `config.go` | modified | Added `metric_categories` + validation |
| `README.md` | modified | Documented the new field and pattern |

---

## 3. Configuration reference

### 3.1 All parameters

| Parameter | Type | Default | Description |
|---|---|---|---|
| `alert_definitions_configmap_name` | string | `opsramp-alert-user-config` | ConfigMap holding alert definitions |
| `alert_definitions_key` | string | `alert-definitions.yaml` | Key within the ConfigMap |
| `namespace` | string | `$NAMESPACE`, else `opsramp-agent` | Namespace of the ConfigMap. Auto-populated; do not set manually |
| `alert_definitions_file_path` | string | — | Absolute path to a `.yaml`/`.yml` file. **Takes precedence over ConfigMap mode** |
| `watch_file_changes` | bool | `true` | Enable file watching (file mode only) |
| `file_watch_interval` | duration string | `30s` | Polling interval when fsnotify is unavailable. Must be positive |
| `metric_categories` | []string | empty (= all) | **New.** Restrict allow-list to `podMetric` and/or `clusterMetric` |

### 3.2 Two mutually exclusive source modes

**ConfigMap mode** (default) builds an in-cluster Kubernetes client and watches
the ConfigMap for changes.

**File mode** activates whenever `alert_definitions_file_path` is set. `Validate()`
then blanks `alert_definitions_configmap_name`, `alert_definitions_key`, and
`namespace` so the two modes cannot conflict. The path must be absolute, must end
in `.yaml`/`.yml`, and must exist at startup.

---

## 4. Metric categories

### 4.1 Categorization rules

| `resourceType` in alert definition | Category assigned |
|---|---|
| `k8s_pod` (also `Pod`, `K8S_POD`, `k8s-pod`) | `podMetric` |
| Any other value (`k8s_cluster`, `k8s_node`, `k8s_deployment`, …) | `clusterMetric` |
| Absent (flat VM-agent format) | `clusterMetric` |

Matching is case-insensitive and treats `-` and `_` as equivalent, because both
`"Pod"` and `"k8s_pod"` spellings appear across OpsRamp payloads.

The category applies to **every metric extracted from that group's expressions**,
since a single `expr` may reference multiple metrics.

### 4.2 Dual-category metrics

A metric referenced by both a `k8s_pod` expression and a non-pod expression
belongs to **both** categories. Categories are accumulated with a bitwise OR
across all rules before any instance reads the result, so rule ordering is
irrelevant:

```go
metrics[name] |= cr.category
```

Worked example:

```yaml
alertDefinitions:
  - resourceType: "k8s_pod"
    rules:
      - name: "pod cpu"
        expr: "container_cpu_usage > 80"
      - name: "pod memory"
        expr: "pod_memory_working_set > 90"
  - resourceType: "k8s_node"
    rules:
      - name: "node cpu"
        expr: "container_cpu_usage / node_cpu_total > 0.9"
```

Resulting categorization:

| Metric | `podMetric` | `clusterMetric` |
|---|---|---|
| `container_cpu_usage` | ✅ | ✅ |
| `pod_memory_working_set` | ✅ | — |
| `node_cpu_total` | — | ✅ |

An instance with `metric_categories: ["podMetric"]` keeps
`container_cpu_usage` + `pod_memory_working_set`; an instance with
`["clusterMetric"]` keeps `container_cpu_usage` + `node_cpu_total`.

### 4.3 Accepted values

| Value | Meaning |
|---|---|
| `podMetric` | Metrics from `k8s_pod` alert definitions |
| `clusterMetric` | Metrics from all other alert definitions |
| *(omitted / empty list)* | All categories — identical to pre-change behavior |

Values are case- and whitespace-insensitive. Unknown values are rejected at
config load with an explicit error listing the valid options.

---

## 5. Usage patterns

### 5.1 Minimal — ConfigMap mode, no categories

Unchanged from before; the allow-list is the union of every alert definition.

```yaml
processors:
  opsrampmetricsfilter:

service:
  pipelines:
    metrics:
      receivers: [otlp]
      processors: [opsrampmetricsfilter, batch]
      exporters: [prometheusremotewrite]
```

### 5.2 Explicit ConfigMap

```yaml
processors:
  opsrampmetricsfilter:
    alert_definitions_configmap_name: "opsramp-alert-user-config"
    alert_definitions_key: "alert-definitions.yaml"
    # namespace comes from the NAMESPACE env var
```

### 5.3 File mode

Useful outside Kubernetes, or when definitions arrive via a mounted volume.

```yaml
processors:
  opsrampmetricsfilter:
    alert_definitions_file_path: "/etc/opsramp/alert-definitions.yaml"
    watch_file_changes: true
    file_watch_interval: "30s"
```

The watcher monitors the **parent directory**, not the file inode, so
ConfigMap volume updates (which swap a symlink rather than writing in place) are
detected. fsnotify events are debounced by 200 ms; if fsnotify is unavailable the
watcher falls back to polling at `file_watch_interval`.

### 5.4 Category split — pod and cluster to different exporters

The primary new pattern. A union pre-filter reduces volume before the fan-out,
then a `forward` connector feeds two category-scoped instances.

```yaml
connectors:
  forward/metricsfilter: {}

processors:
  k8s_attributes/worker:
  batch/opsramp-metrics:

  opsrampmetricsfilter:                    # union pre-filter

  opsrampmetricsfilter/pod:
    metric_categories: ["podMetric"]

  opsrampmetricsfilter/cluster:
    metric_categories: ["clusterMetric"]

service:
  pipelines:
    metrics/controlplane-node:
      receivers: [kubelet_stats, prometheus/kube-proxy]
      processors: [k8s_attributes/worker, opsrampmetricsfilter]
      exporters: [opsrampdebug, forward/metricsfilter]

    metrics/pod:
      receivers: [forward/metricsfilter]
      processors: [opsrampmetricsfilter/pod, batch/opsramp-metrics]
      exporters: [exporter1]

    metrics/cluster:
      receivers: [forward/metricsfilter]
      processors: [opsrampmetricsfilter/cluster, batch/opsramp-metrics]
      exporters: [exporter2]
```

Data flow:

```
kubelet_stats ─┐
               ├─▶ k8s_attributes ─▶ union filter ─┬─▶ opsrampdebug
kube-proxy ────┘                                   │
                                                   └─▶ forward ─┬─▶ pod filter ─▶ exporter1
                                                                └─▶ cluster filter ─▶ exporter2
```

Notes on this topology:

- **One connector, N pipelines.** A `forward` connector may be listed as the
  receiver of any number of pipelines; it fans out like a receiver. Adding a
  branch means adding a *pipeline*, not another connector.
- **Keep the union pre-filter.** The fan-out deep-copies the batch for each
  mutating consumer, so filtering first is what keeps the copy cost small.
- **Batch downstream of the split.** Batching before the fan-out just gets torn
  apart again. Note that reusing the same processor ID across pipelines creates a
  *separate instance* per pipeline, so listing `batch` in both the upstream and
  downstream pipelines pays the timeout latency twice.
- Metrics in both categories are emitted by both branches — no extra config.

### 5.5 Overlapping categories

An instance may select multiple categories. Useful when one exporter should
receive everything while another receives only pod metrics.

```yaml
processors:
  opsrampmetricsfilter/all:
    metric_categories: ["podMetric", "clusterMetric"]   # same as omitting the field
  opsrampmetricsfilter/pod:
    metric_categories: ["podMetric"]
```

### 5.6 Distinct sources

Instances pointing at different ConfigMaps each get their own loader — sharing is
keyed on the source, not on the component name.

```yaml
processors:
  opsrampmetricsfilter:
    alert_definitions_configmap_name: "test-alert-config"
    namespace: "test-namespace"
  opsrampmetricsfilter/custom:
    alert_definitions_configmap_name: "custom-alerts"
    alert_definitions_key: "custom-alerts.yaml"
    namespace: "monitoring"
```

---

## 6. Architecture

### 6.1 Loader sharing key

Instances share a loader when this key matches:

- **File mode:** `alert_definitions_file_path`
- **ConfigMap mode:** `namespace` + `alert_definitions_configmap_name` + `alert_definitions_key`

`metric_categories` is deliberately **not** part of the key — that is what allows
pod and cluster instances to share one watch.

### 6.2 Lifecycle

1. `newFilterProcessor` parses `metric_categories` into a bitmask.
2. `acquireLoader` returns the existing loader for the source or creates one,
   incrementing the refcount.
3. The instance subscribes; the current snapshot is delivered immediately, so a
   late-joining instance need not wait for the next reload.
4. `loader.start()` (idempotent) performs the initial load and begins watching.
5. On each reload the loader publishes a new snapshot to every subscriber.
6. `Shutdown` unsubscribes and releases; the last release cancels the context,
   closes the file watcher, and removes the loader from the registry.

### 6.3 Data flow through a reload

```
ConfigMap change
   └─▶ loader.reload()
        └─▶ parse YAML (nested, else flat fallback)
             └─▶ per rule: PromQL AST walk → metric names
                  └─▶ metrics[name] |= category
                       └─▶ publish(map[string]categorySet)
                            ├─▶ pod instance:     keep where set & categoryPod
                            └─▶ cluster instance: keep where set & categoryCluster
```

The published map is shared read-only across subscribers and is never mutated
after publish. Each instance projects it into its own `map[string]bool` and swaps
that map in under a write lock, so readers never observe a partially built map.

### 6.4 Filtering behavior (unchanged)

- Metric names are normalized `.` → `_` before lookup, so OTel-style
  `k8s.pod.cpu.usage` matches PromQL-style `k8s_pod_cpu_usage`.
- Resource and scope containers are created lazily, so empty ones are not emitted.
- Empty incoming batch → passed through untouched.
- **Empty allow-list → all metrics dropped.**

---

## 7. Alert definitions formats

### 7.1 Nested (preferred)

Carries `resourceType`, so categorization works.

```yaml
alertDefinitions:
  - resourceType: "k8s_cluster"
    rules:
      - name: "k8s_apiserver_requests_error_rate"
        interval: "30s"
        expr: "sum(rate(apiserver_request_total[5m])) > 10"
        isAvailability: false
        warnOperator: ">"
        warnThreshold: "5"
  - resourceType: "k8s_pod"
    rules:
      - name: "k8s_pod_cpu_usage"
        expr: "k8s_pod_cpu_limit_utilization_ratio * 100"
```

### 7.2 Flat (VM-agent fallback)

Used automatically when the nested parse yields no rules. It carries no
`resourceType`, so **all** rules are treated as `clusterMetric`.

```yaml
alertDefinitions:
  - name: "flat rule"
    expr: "some_metric > 1"
```

> A `podMetric`-scoped instance reading flat-format definitions produces an empty
> allow-list and therefore drops everything.

---

## 8. Backward compatibility

No configuration changes are required. All six pre-existing parameters are
unchanged, defaults are unchanged, and omitting `metric_categories` reproduces
the previous allow-list exactly.

### Behavioral deltas

| Change | Impact |
|---|---|
| "Extracted metrics from alert definitions" (full metric list per reload) removed | Replaced by "Loaded alert definitions" / "Published alert definitions snapshot" / "Applied alert definitions" |
| Loader-level logs use the first-acquiring instance's logger | With a category split, shared loader logs are attributed to whichever instance was constructed first |
| ConfigMap `Get`/`Watch` use the loader context instead of `context.TODO()` | Cancellable on shutdown; an in-flight `Get` during shutdown now logs a context-cancelled error |
| Empty metric names skipped during extraction | `{__name__=~"..."}` selectors yield `""`, which could never match a real metric |
| `watch_file_changes` / `file_watch_interval` are first-acquire-wins | Only affects multiple instances sharing one source with *different* watcher settings |

`ConsumeMetrics` is unchanged — the filtering logic and all of its log statements
(including per-batch `Info` logging) are byte-identical to the previous version.

---

## 9. Known limitations

**Empty allow-list drops everything, quietly.** If the ConfigMap read fails, the
YAML is unparseable, the ConfigMap is deleted, or a selected category matches no
rules, the instance drops 100% of metrics and logs only at `Info`. With a
category split this is one silent-blackhole risk per branch. Recommended
mitigation: alert on the `"No metrics filter configured, dropping all metrics"`
log line.

**A valid-but-unmatched category cannot be distinguished from "no rules yet."**
`metric_categories` values are validated against the known set, but a
syntactically valid category that matches zero alert definitions is
indistinguishable from a source that legitimately has none.

**`ConsumeMetrics` holds a read lock across the downstream call.** A slow
exporter delays reload application for that instance. Pre-existing behavior,
retained.

**Fan-out deep-copies per mutating consumer.** The processor declares
`MutatesData: true`, so an N-way split costs N−1 full batch clones. Mitigate by
placing a union pre-filter upstream of the fan-out (see §5.4).

---

## 10. Alternative considered: marker attributes + routing connector

An earlier design had the processor stamp `opsramp.pod_metric` /
`opsramp.cluster_metric` datapoint attributes and split with `routingconnector`
using `action: copy`.

Rejected because datapoint attributes become **labels** in
`prometheusremotewrite`, changing series identity, breaking continuity with
previously written data, and inflating cardinality. Avoiding it would require a
`transform` processor in every downstream pipeline to strip the markers.

The category-scoped instance approach needs no marker attributes, no OTTL, and no
strip step. Its cost is the extra batch clones noted above, which the union
pre-filter largely offsets.
