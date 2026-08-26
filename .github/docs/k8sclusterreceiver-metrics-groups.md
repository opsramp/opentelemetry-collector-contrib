# k8sclusterreceiver — `metrics_groups` enhancement

Per-resource-type collection control for the `k8s_cluster` receiver.

- **Component:** `receiver/k8sclusterreceiver`
- **Branch:** `release/v0.159.0` (OpsRamp fork, not submitted upstream)
- **Commit:** `0922a9a572a`

---

## 1. Problem

The receiver exposed only the mdatagen-generated `metrics:` block, which toggles individual metrics
on and off. That control acts at the *emission* layer: the metric is computed and then discarded.

Every Kubernetes resource type is watched regardless. On startup the receiver opens **13
unconditional LIST/WATCH streams** against the API server and keeps a full client-side cache for
each. Three more (`EndpointSlice`, `PersistentVolume`, `PersistentVolumeClaim`) are conditional.

The consequence: an operator who only wants node and deployment metrics still pays for a Pod
informer holding every pod in the cluster. In large clusters the Pod cache dominates the
collector's resident memory, and the watch stream dominates its API server load. No amount of
`metrics:` tuning could reduce it, and neither could a downstream `filterprocessor`.

## 2. Solution

A new `metrics_groups` config block that operates at the **informer** layer. Each group maps 1:1 to
a Kubernetes kind. Disabling a group means the informer is never created, so:

- no LIST/WATCH is issued for that resource,
- no client-side cache is allocated,
- no metrics are emitted,
- the corresponding RBAC `list`/`watch` permission is no longer required.

```yaml
k8s_cluster:
  metrics_groups:
    pods:
      enabled: false
    persistentvolumeclaims:
      enabled: false
```

An `enabled_by_default` flag supports the opt-out style, so a narrow configuration does not have to
enumerate all 16 groups:

```yaml
k8s_cluster:
  metrics_groups:
    enabled_by_default: false
    nodes:
      enabled: true
    deployments:
      enabled: true
```

## 3. Design

### 3.1 Why the change is small

`DataCollector.CollectMetricData` was already purely store-driven — it iterates
`metadataStore.ForEach(gvk.X, ...)` for each resource. `Store.ForEach` ranges over
`ms.stores[gvk]`, and a missing key yields a nil map, which ranges zero times.

Skipping `setupInformer` therefore makes the metrics disappear with **no collector-side change and
no nil-pointer path**. The entire production change is 12 lines in `watcher.go` plus the config
types. No generated file, metric definition, or function signature was touched.

### 3.2 Group-to-kind mapping

| Group | Kind | Metrics |
| --- | --- | --- |
| `pods` | `Pod` | `k8s.pod.*` **and all `k8s.container.*`** |
| `nodes` | `Node` | `k8s.node.*`, incl. `node_conditions_to_report` / `allocatable_types_to_report` |
| `namespaces` | `Namespace` | `k8s.namespace.phase` |
| `deployments` | `Deployment` | `k8s.deployment.*` |
| `replicasets` | `ReplicaSet` | `k8s.replicaset.*` |
| `replicationcontrollers` | `ReplicationController` | `k8s.replication_controller.*` |
| `daemonsets` | `DaemonSet` | `k8s.daemonset.*` |
| `statefulsets` | `StatefulSet` | `k8s.statefulset.*` |
| `jobs` | `Job` | `k8s.job.*` |
| `cronjobs` | `CronJob` | `k8s.cronjob.*` |
| `horizontalpodautoscalers` | `HorizontalPodAutoscaler` | `k8s.hpa.*` |
| `resourcequotas` | `ResourceQuota` | `k8s.resource_quota.*` |
| `services` | `Service` + `EndpointSlice` | `k8s.service.*` |
| `persistentvolumes` | `PersistentVolume` | `k8s.persistentvolume.*` |
| `persistentvolumeclaims` | `PersistentVolumeClaim` | `k8s.persistentvolumeclaim.*` |
| `clusterresourcequotas` | `ClusterResourceQuota` | `openshift.*quota.*` (openshift only) |

Two deliberate groupings:

- **No `containers` group.** Container metrics are recorded inside `pod.RecordMetrics`, which loops
  over `pod.Spec.Containers`. There is no Container informer. A separate group would be misleading
  because the expensive part — the Pod watch — would still run.
- **`EndpointSlice` folded into `services`.** It emits no metrics of its own; it exists only to back
  `k8s.service.endpoint.count`.

### 3.3 Backward compatibility via `*bool`

Both `ResourceGroupConfig.Enabled` and `MetricsGroupsConfig.EnabledByDefault` are `*bool`, so
"unset" is representable and distinct from `false`:

```go
func (c MetricsGroupsConfig) enabledForKind(kind string) bool {
	if enabled := c.groupForKind(kind).Enabled; enabled != nil {
		return *enabled
	}
	return c.EnabledByDefault == nil || *c.EnabledByDefault
}
```

Defaults were deliberately **not** populated in `createDefaultConfig()`. Had a plain `bool` been
used, anything constructing `&Config{...}` directly — which the existing tests and any future
upstream code do — would have silently disabled every group. The nil-pointer approach makes the
safe path the default path.

Unknown kinds fall through to `enabled_by_default`, so a kind added by a future upstream merge is
not silently dropped when the operator has not opted out.

### 3.4 Gating logic

```go
// watcher.go — after the existing supportedKinds construction
for kind := range supportedKinds {
	if !rw.config.MetricsGroups.enabledForKind(kind) {
		delete(supportedKinds, kind)
		rw.logger.Info("Metrics group is disabled, resource will not be watched",
			zap.String("kind", kind))
	}
}
```

The OpenShift quota client creation in `initialize()` is gated the same way, so
`clusterresourcequotas: {enabled: false}` avoids building that client entirely.

Deleting from a map during `range` is legal in Go; deleted keys are simply not produced.

### 3.5 Validation

`Validate()` rejects a config where every group ends up disabled, which would otherwise produce a
receiver that silently collects nothing:

```
all metrics_groups are disabled, the receiver would not collect anything
```

## 4. Files changed

| File | Change |
| --- | --- |
| `config.go` | `MetricsGroupsConfig`, `ResourceGroupConfig`, `groupForKind`, `enabledForKind`, `anyEnabled`, validation |
| `watcher.go` | Informer gating loop; OpenShift quota client guard |
| `config_test.go` | Group resolution matrix, opt-out mode, validation, empty-block case |
| `watcher_test.go` | Informer-set assertions incl. backward-compat pin |
| `testdata/config.yaml` | `metrics_groups` and `empty_metrics_groups` fixtures |
| `README.md` | Group table, `enabled_by_default`, combining with `metrics`, metadata impact |

---

## 5. Performance

### 5.1 Method

All figures below are **measured**, produced by benchmarks added alongside the change:

- `internal/collection/collector_benchmark_test.go` — `BenchmarkCollectMetricData`
- `metrics_groups_benchmark_test.go` — `BenchmarkInformerCacheObject`,
  `BenchmarkPrepareSharedInformerFactory`, `BenchmarkInformerCacheGrowth`

Command:

```bash
cd receiver/k8sclusterreceiver
go test -run XXX -bench . -benchmem -benchtime 3s ./...
```

Environment:

| | |
| --- | --- |
| Go | 1.26.4 |
| Platform | darwin/arm64 |
| CPU | Apple M1 Max (10 threads) |
| macOS | 26.6.2 |

Workload model: pod counts of 500 / 2000 / 5000, each pod carrying one container, plus 50 objects
each of Node, Namespace, Deployment, ReplicaSet, DaemonSet, StatefulSet, Job, CronJob, HPA,
ResourceQuota, ReplicationController and Service. `pods_disabled` registers no Pod store at all,
which is exactly the state produced by `metrics_groups: {pods: {enabled: false}}`.

> **Scope caveat.** These measure in-process cost: scrape CPU, allocation, and cached-object
> footprint. They do **not** measure API server load, network traffic, or steady-state RSS on a live
> cluster — that requires a real deployment. Watch-stream counts in §5.5 are structural facts read
> from the code, not measurements.

### 5.2 Per-scrape cost — `BenchmarkCollectMetricData`

| Scenario | Time/scrape | Allocated/scrape | Allocations/scrape |
| --- | --- | --- | --- |
| `pods_disabled` | **1.63 ms** | **1.01 MB** | **21,992** |
| 500 pods | 5.41 ms | 3.09 MB | 64,504 |
| 2,000 pods | 16.24 ms | 9.26 MB | 192,015 |
| 5,000 pods | 40.21 ms | 21.46 MB | 447,017 |

Savings from `pods: {enabled: false}`:

| Cluster size | CPU reduction | Allocation reduction |
| --- | --- | --- |
| 500 pods | 70 % | 67 % |
| 2,000 pods | 90 % | 89 % |
| 5,000 pods | **96 %** | **95 %** |

Marginal cost of the pods group, derived from the 5,000-pod row:

- **7.7 µs** of CPU per pod per scrape
- **4.09 KB** allocated per pod per scrape

At the default 10 s collection interval, a 5,000-pod cluster spends ~40 ms of CPU every 10 s
(~0.4 % of one core) and churns ~21 MB of garbage per scrape, ~129 MB/min. Disabling the pods group
removes ~95 % of that.

Scaling is linear in pod count, as expected from the `ForEach` loop:

```
500 → 2000 (4.0×):   5.41 → 16.24 ms   (3.0× on the pod-attributable portion)
2000 → 5000 (2.5×): 16.24 → 40.21 ms   (2.6× on the pod-attributable portion)
```

### 5.3 Cached-object footprint — `BenchmarkInformerCacheObject`

Bytes retained per object after the informer transform, i.e. the per-object price of keeping an
informer running:

| Kind | Bytes/object | Allocations | Time |
| --- | --- | --- | --- |
| Node | 2,704 B | 5 | 904.0 ns |
| Pod | 1,904 B | 3 | 894.9 ns |
| Deployment | 1,280 B | 1 | 357.7 ns |
| Job | 1,408 B | 1 | 462.3 ns |
| ReplicaSet | 1,152 B | 1 | 362.5 ns |
| Service | 640 B | 1 | 194.7 ns |

Note this is the **post-transform** size. `transformObject` already strips unused fields; the raw
API objects are substantially larger, so these are lower bounds on real cache cost.

### 5.4 Cache growth — `BenchmarkInformerCacheGrowth`

| Pods cached | Bytes | Per pod |
| --- | --- | --- |
| 500 | 0.96 MB | 1,920 B |
| 2,000 | 3.84 MB | 1,920 B |
| 5,000 | 9.60 MB | 1,920 B |

Perfectly linear at **1,920 B/pod**. A 5,000-pod cluster holds ~9.6 MB of pod objects, reclaimed
entirely by disabling the group. Real-world figures will be higher: this excludes client-go's
`ThreadSafeStore` index maps, key strings, and the DeltaFIFO used during resync.

### 5.5 Startup — `BenchmarkPrepareSharedInformerFactory`

| Configuration | Time | Allocated | Allocations | Watch streams |
| --- | --- | --- | --- | --- |
| All groups (default) | 95,902 ns | 309,936 B | 611 | 13 |
| `nodes` + `deployments` only | **31,285 ns** | **58,364 B** | **156** | **2** |
| *Improvement* | *67 % faster* | *81 % less* | *74 % fewer* | *85 % fewer* |

The watch-stream column is the operationally significant one. Each eliminated informer removes a
persistent HTTP/2 stream to the API server, its initial full LIST, and its resync traffic every
`metadata_collection_interval` (default 5 m).

With `namespaces: [ns1, ns2, ...]` set, namespaced kinds get **one informer per namespace**, so the
watch-count saving multiplies by the namespace count.

### 5.6 Summary

For a 5,000-pod cluster disabling the pods group:

| Dimension | Before | After | Change |
| --- | --- | --- | --- |
| CPU per scrape | 40.21 ms | 1.63 ms | −96 % |
| Allocation per scrape | 21.46 MB | 1.01 MB | −95 % |
| Garbage per minute @10 s | ~129 MB | ~6 MB | −95 % |
| Pod cache retained | ~9.6 MB | 0 MB | −100 % |
| Watch streams | 13 | 12 | −1 |

The `metrics:` toggles achieve **none** of the above — they suppress emission only, after the data
has been fetched, cached and computed.

---

## 6. Operational notes

### 6.1 Metadata and entity events are coupled to the informer

`setupInformer` does two jobs: it registers the cache used for metrics **and** attaches
`AddFunc`/`UpdateFunc`/`DeleteFunc` handlers that drive metadata updates and entity events.
Disabling a group stops both. Affected destinations:

- **`metadata_exporters`** — only the signalfx exporter implements `MetadataExporter` in contrib. It
  converts labels/annotations into dimension property updates. When a group is disabled those
  updates stop and previously reported properties go **stale** rather than disappearing.
- **Entity events** (`otel.entity.*` log records) — emitted only when the receiver sits in a `logs`
  pipeline. `metadata_collection_interval` is the informer resync period, so removing the informer
  also stops periodic re-emission of entity state.

This intentionally overrides `shouldWatchResourceForMetadataOnly()`, which would otherwise keep the
PV/PVC informers alive whenever a metadata destination is configured.

### 6.2 Cross-resource metadata dependency

Pod metadata enrichment reads three other caches:

```go
if store := mc.Get(gvk.Service);    store != nil { ...GetPodServiceTags... }
if store := mc.Get(gvk.Job);        store != nil { ...collectPodJobProperties... }
if store := mc.Get(gvk.ReplicaSet); store != nil { ...collectPodReplicaSetProperties... }
```

Disabling `services`, `jobs` or `replicasets` while `pods` stays enabled drops service tags and the
resolved owning workload (ReplicaSet→Deployment, Job→CronJob) from pod metadata. It **degrades
safely** — the nil guards and `getObjectFromStore` returning `(nil, nil)` mean it only logs at debug
level. **Pod metrics are unaffected.**

### 6.3 Precedence with `metrics:`

| `metrics_groups` | `metrics` | Result |
| --- | --- | --- |
| enabled | enabled | emitted |
| enabled | disabled | not emitted, informer still runs |
| disabled | enabled | **not emitted — group wins** |
| disabled | disabled | not emitted |

Enabling a group never forces metrics on; it makes the `metrics` toggles meaningful. A metric cannot
be revived under `metrics` once its group is off, because no informer feeds it.

Precedence is documented rather than validated: the receiver cannot distinguish "explicitly enabled"
from "enabled by default" because `enabledSetByUser` is unexported in `internal/metadata`.

### 6.4 Not affected

The `k8sattributes` processor and `kubeletstats` receiver run their own watchers and are unchanged.
Span/log/metric pod enrichment and kubelet-sourced pod metrics keep working normally.

### 6.5 RBAC

Disabled groups no longer need their `list`/`watch` permissions, so the ClusterRole can be trimmed
accordingly.

---

## 7. Verification

| Check | Result |
| --- | --- |
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `gofmt -l .` | clean |
| `go test ./...` | pass — receiver + 20 internal packages |

Backward compatibility is pinned at three levels:

1. **By construction** — all-nil pointers resolve to enabled; the gating loop deletes nothing.
2. `TestNoMetricsGroupsWatchesDefaultKinds` — asserts the informer set from `createDefaultConfig()`
   is exactly the pre-change one (13 kinds present; `EndpointSlice`/PV/PVC absent).
3. `TestLoadConfig/k8s_cluster` and `.../empty_metrics_groups` — bare and null `metrics_groups:`
   both equal `createDefaultConfig()`.

New tests: `TestMetricsGroupsEnabledForKind`, `TestMetricsGroupsDisableInformers`,
`TestMetricsGroupsEnabledByDefaultOptOut`, `TestNoMetricsGroupsWatchesDefaultKinds`, plus validation
and config-loading cases.

---

## 8. Reproducing the benchmarks

```bash
cd receiver/k8sclusterreceiver

# everything
go test -run XXX -bench . -benchmem -benchtime 3s ./...

# scrape cost only
go test -run XXX -bench BenchmarkCollectMetricData -benchmem ./internal/collection/

# startup and cache footprint only
go test -run XXX -bench 'BenchmarkInformerCache|BenchmarkPrepareSharedInformerFactory' -benchmem .
```

For statistically sound comparisons, use `-count=10` and `benchstat`.

## 9. Possible follow-ups

- Measure real RSS and API-server QPS on a live cluster to validate §5.4 and §5.5 end-to-end.
- Optional label/field selectors per group to narrow watches further without disabling them.
- An internal telemetry counter for objects held per informer, to make the cost observable at
  runtime rather than only in benchmarks.
