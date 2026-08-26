# kubeletstatsreceiver — Pod & Container State Metrics Enhancement

> **Scope:** OpsRamp fork enhancement to `receiver/kubeletstatsreceiver`.
> **Commit:** `ae77982` on `release/v0.159.0`.
> **Purpose:** move pod/container state metric collection from the cluster-wide
> `k8sclusterreceiver` to the per-node DaemonSet agent.
>
> Unless stated otherwise, file paths in this document are relative to
> `receiver/kubeletstatsreceiver/`.

---

## 1. Motivation

`k8sclusterreceiver` runs as a single cluster-wide collector and produces pod and container
state metrics (`k8s.pod.phase`, `k8s.container.restarts`, …) by maintaining an informer cache
of **every pod in the cluster**, backed by a watch against the API server.

That design has three costs:

1. **API server load** — a cluster-wide pod watch.
2. **Memory** — the informer cache holds all pods; hundreds of MB on large clusters.
3. **Single point of failure** — one collector owns all pod state.

The kubelet on each node already exposes everything needed, locally:

| Endpoint | Contents |
| --- | --- |
| `/stats/summary` | **Usage** metrics — what kubeletstats already scrapes |
| `/pods` | **Full Pod objects for pods on that node** — `status.phase`, `status.containerStatuses[].restartCount`, `.ready`, `.state.reason`, `spec.containers[].resources.limits/requests` |

`kubeletstatsreceiver` already had the machinery to call `/pods` (used for metadata enrichment
such as `container.id` and volume types), but **never converted that data into metrics**.

This enhancement closes that gap.

### Resulting architecture

```
┌──────────────── DaemonSet agent (per node) ─────────────────┐
│  kubeletstats → local kubelet /stats/summary → USAGE        │
│  kubeletstats → local kubelet /pods          → STATE  [NEW] │
│      (phase, restarts, ready, reason, limits, requests)     │
└─────────────────────────────────────────────────────────────┘
```

No API server involvement. No kube-state-metrics. Cost is bounded by pods-per-node
(default cap 110, max 250) rather than pods-per-cluster.

---

## 2. Metrics Added

Eleven metrics. Names, units, metric types, value types and value encodings are **identical to
`k8sclusterreceiver`** — verified programmatically by diffing both `metadata.yaml` files.

### 2.1 Enabled by default

| Metric | Type | Value | Unit | Source field |
| --- | --- | --- | --- | --- |
| `k8s.pod.phase` | Gauge | Int | — | `pod.Status.Phase` |
| `k8s.pod.status_reason` | Gauge | Int | — | `pod.Status.Reason` |
| `k8s.container.restarts` | Gauge | Int | `{restart}` | `containerStatus.RestartCount` |
| `k8s.container.ready` | Gauge | Int | — | `containerStatus.Ready` |
| `k8s.container.status.reason` | Sum (cumulative, non-monotonic) | Int | `{container}` | `containerStatus.State.Waiting/Terminated.Reason` |

### 2.2 Disabled by default (opt-in)

| Metric | Type | Value | Unit | Source field |
| --- | --- | --- | --- | --- |
| `k8s.container.cpu_limit` | Gauge | Double | `{cpu}` | `spec.containers[].resources.limits.cpu` |
| `k8s.container.cpu_request` | Gauge | Double | `{cpu}` | `…requests.cpu` |
| `k8s.container.memory_limit` | Gauge | Int | `By` | `…limits.memory` |
| `k8s.container.memory_request` | Gauge | Int | `By` | `…requests.memory` |
| `k8s.container.storage_limit` | Gauge | Int | `By` | `…limits.storage` |
| `k8s.container.storage_request` | Gauge | Int | `By` | `…requests.storage` |

### 2.3 Value encodings

`k8s.pod.phase` — single data point, integer-encoded:

| Value | Phase |
| --- | --- |
| 1 | Pending |
| 2 | Running |
| 3 | Succeeded |
| 4 | Failed |
| 5 | Unknown *(also the fallback for unrecognised values)* |

`k8s.pod.status_reason` — single data point, integer-encoded:

| Value | Reason |
| --- | --- |
| 1 | Evicted |
| 2 | NodeAffinity |
| 3 | NodeLost |
| 4 | Shutdown |
| 5 | UnexpectedAdmissionError |
| 6 | Unknown *(also the fallback for an empty reason)* |

`k8s.container.ready` — `0` = not ready, `1` = ready.

`k8s.container.restarts` — raw `RestartCount`. Note the kubelet may prune dead containers and
reset this to 0, so treat it as `== 0` vs `> 0` rather than trusting the absolute value.

`k8s.container.cpu_*` — computed as `MilliValue() / 1000.0`, matching k8scluster exactly
(not `AsApproximateFloat64()`).

`k8s.container.{memory,storage}_*` — `Quantity.Value()` in bytes.

Limits and requests are only emitted when actually present in the pod spec; an absent limit
produces no data point rather than a zero.

### 2.4 `k8s.container.status.reason` — the 9-data-point fan-out

This metric emits **one data point per known reason on every scrape**, with `1` for the current
reason and `0` for all others:

```
ContainerCreating, CrashLoopBackOff, CreateContainerConfigError, ErrImagePull,
ImagePullBackOff, OOMKilled, Completed, Error, ContainerCannotRun
```

**This is deliberate and matches upstream k8scluster** (upstream commit `aaf8a0d3e122`,
2025-09-14). The rationale, from the upstream metric description:

> All possible container state reasons will be reported at each time interval **to avoid missing
> metrics**. Only the value corresponding to the current state reason will be non-zero.

If instead only the *current* reason were emitted, a container transitioning
`CrashLoopBackOff → Running` would cause the `CrashLoopBackOff` series to simply stop being
reported — it would go **stale**, not to zero. Prometheus-style backends keep serving the last
value through the staleness window, so `sum by (reason)` queries silently over-count.

**Cost:** this single metric accounts for ~40% of the added data point volume (see §6).
Upstream ships it `enabled: false`; this fork ships it `enabled: true` by explicit decision.
Set `k8s.container.status.reason: {enabled: false}` to reclaim that volume.

---

## 3. Implementation

### 3.1 Files changed

| File | Change |
| --- | --- |
| `metadata.yaml` | 11 metric definitions + `k8s.container.status.reason` attribute enum |
| `internal/kubelet/podstatus.go` | **New** — all state metric emission logic |
| `internal/kubelet/podstatus_test.go` | **New** — unit tests |
| `internal/kubelet/metrics.go` | Pod-list iteration pass; new `podStatusEnabled` parameter |
| `scraper.go` | `needsPodStatus` flag gating the `/pods` fetch |
| `internal/metadata/generated_*.go` | mdatagen output |
| `documentation.md` | mdatagen output |
| `README.md` | New "Pod and container state metrics" section + RBAC updates |
| `testdata/pods.json` | Added missing `spec.containers` (see §3.5) |
| `testdata/scraper/*.yaml` | 13 regenerated golden files |

Plus `.chloggen/kubeletstats-pod-status-metrics.yaml` at the repository root.

### 3.2 Fetch gating

`/pods` is **not** unconditionally fetched. It was already conditional, and the new flag simply
joins the existing condition — there is at most **one** `/pods` call per scrape, never two:

```go
// scraper.go
if len(r.extraMetadataLabels) > 0 || r.needsResources || r.needsPodStatus {
    podsMetadata, err = r.metadataProvider.Pods()
    ...
}
```

`needsPodStatus` is true if **any** of the 11 metrics is enabled. Since five are on by default,
`/pods` is now fetched by default — this is the main operational consequence (see §7).

### 3.3 Driven by the pod list, not the stats summary

State metrics are emitted from `Metadata.PodsMetadata` (i.e. `/pods`) **after** the usage pass,
rather than being folded into `podStats` / `containerStats`:

```go
// internal/kubelet/metrics.go
acc.nodeStats(summary.Node)
for i := range summary.Pods { /* usage metrics */ }

// Pod state metrics come from the /pods endpoint rather than the stats summary so that pods and
// containers without usage stats are still reported.
if podStatusEnabled && acc.metadata.PodsMetadata != nil {
    for i := range acc.metadata.PodsMetadata.Items {
        acc.podStatus(&acc.metadata.PodsMetadata.Items[i])
    }
}
```

**Why this matters:** `/stats/summary` only contains pods/containers that have running cgroups.
A pod stuck in `Pending`, `ImagePullBackOff`, or a container in `CrashLoopBackOff` may be absent
or partially absent from the summary — yet those are precisely the states these metrics exist to
detect. Driving off the pod list guarantees complete coverage.

**Trade-off:** each pod and container gets a *second* `ResourceMetrics` group with identical
resource attributes to its usage counterpart (see §6.2).

### 3.4 Resource attributes

State metrics carry the **same resource attributes as the corresponding kubeletstats usage
metrics**, so the two sets share a resource identity and merge downstream:

| Level | Attributes |
| --- | --- |
| Pod | `k8s.pod.uid`, `k8s.pod.name`, `k8s.namespace.name` |
| Container | above + `k8s.container.name` (+ `container.id` only when `extra_metadata_labels: [container.id]`) |

Metric group gating is respected: pod state metrics require the `pod` metric group, container
state metrics require the `container` group.

> **Known divergence from k8scluster** — see §5.2. `k8s.node.name` was deliberately *not* added.

### 3.5 Test fixture correction

`testdata/pods.json` contained pods with `status.containerStatuses` but **no `spec.containers`** —
a shape a real kubelet never returns. With the state metrics off, nothing exercised that path;
once enabled, only 1 of 9 containers produced state metrics.

Rather than adding defensive code for a case that cannot occur in production, the **fixture was
corrected**: 8 name-only `spec.containers` entries were added to match the existing container
statuses.

Verified this does not perturb existing behaviour: a container with no `resources` block yields
the same zero-valued `podResources` aggregation as no container at all, so the
`*_utilization` metric tests are unaffected.

---

## 4. Testing

| Test | Coverage |
| --- | --- |
| `TestPodStatusMetrics` | All metrics, values, resource attributes, `container.id`, reason fan-out; includes a pod **absent from the stats summary** |
| `TestPodStatusMetricsDisabled` | No emission when disabled |
| `TestPodStatusMetricGroupsFiltered` | `metric_groups` gating |
| `TestPodPhaseToInt` / `TestPodStatusReasonToInt` | Full enum mapping incl. fallbacks |
| `TestScraperWithPodStatusMetrics` | End-to-end through the scraper, incl. the `/pods` fetch |
| 13 golden-file tests | Regenerated and passing |

Metric-count constants updated in `scraper_test.go`:

| Constant | Before | After | Delta |
| --- | --- | --- | --- |
| `nodeMetrics` | 15 | 15 | — |
| `podMetrics` | 15 | 17 | `+phase +status_reason` |
| `containerMetrics` | 11 | 22 | `+restarts +ready +9 status.reason` |
| `volumeMetrics` | 5 | 5 | — |

`gofmt`, `go vet` and the full test suite pass.

---

## 5. Parity with k8sclusterreceiver

### 5.1 Metric shape — full parity ✅

Verified by diffing both `metadata.yaml` files on unit, metric kind, value type, aggregation
temporality, monotonicity and attributes:

```
k8s.pod.phase                    MATCH   gauge                      int     unit=""
k8s.pod.status_reason            MATCH   gauge                      int     unit=""
k8s.container.restarts           MATCH   gauge                      int     unit="{restart}"
k8s.container.ready              MATCH   gauge                      int     unit=""
k8s.container.status.reason      MATCH   sum(mono=false,cumulative) int     unit="{container}"
                                                                     attrs=[k8s.container.status.reason]
k8s.container.cpu_limit          MATCH   gauge                      double  unit="{cpu}"
k8s.container.cpu_request        MATCH   gauge                      double  unit="{cpu}"
k8s.container.memory_limit       MATCH   gauge                      int     unit="By"
k8s.container.memory_request     MATCH   gauge                      int     unit="By"
k8s.container.storage_limit      MATCH   gauge                      int     unit="By"
k8s.container.storage_request    MATCH   gauge                      int     unit="By"
```

### 5.2 Resource attributes — intentional divergence ⚠️

| Level | k8sclusterreceiver | kubeletstatsreceiver (this change) |
| --- | --- | --- |
| Pod | `k8s.pod.uid`, `k8s.pod.name`, `k8s.namespace.name`, **`k8s.node.name`**, **`k8s.pod.qos_class`** | `k8s.pod.uid`, `k8s.pod.name`, `k8s.namespace.name` |
| Container | above + `k8s.container.name`, **`container.id`**, **`container.image.name`**, **`container.image.tag`**, **`k8s.container.status.last_terminated_reason`** | `k8s.container.name` + pod attrs (+ `container.id` opt-in) |

**Rationale:** kubeletstats' *existing* usage metrics use the narrower set. Matching them keeps
state and usage metrics on one resource identity. Adopting k8scluster's wider set would
permanently split them into parallel `ResourceMetrics` streams.

**Migration impact:** queries grouping or filtering by `k8s.node.name`,
`container.image.name`, `container.image.tag`, `k8s.pod.qos_class` or
`k8s.container.status.last_terminated_reason` **will break**. Queries keyed on pod name,
pod UID, namespace or container name port over unchanged.

> Adding `k8s.node.name` to all pod/container resources was implemented and then **reverted** by
> decision, as it changes existing usage-metric output. If needed, prefer adding it in the
> pipeline via `resourcedetection` or `k8sattributes` — in a DaemonSet the node is fixed anyway.

### 5.3 Not ported

`k8s.container.status.state` (3-way `running`/`waiting`/`terminated` fan-out) exists in
k8scluster but was not requested and is not implemented here.

---

## 6. Performance

### 6.1 Methodology

Go benchmarks against a synthetic-but-realistic node, then discarded (not committed).

| Parameter | Value |
| --- | --- |
| Hardware | Apple M1 Max, `darwin/arm64` |
| Node size | 60 pods × 2 containers = **120 containers** |
| `/pods` payload | **352,135 bytes** (realistic: labels, annotations, env vars, probes, volumes, tolerations, conditions, container statuses) |

### 6.2 CPU and allocations, per scrape

| Step | Time | Allocated | Allocations |
| --- | ---: | ---: | ---: |
| `/pods` fetch + JSON unmarshal | 4.47 ms | 1.25 MB | 17,298 |
| `MetricsData` — *without* state metrics | 0.49 ms | 555 KB | 9,536 |
| `MetricsData` — *with* state metrics | 1.03 ms | 1.19 MB | 23,937 |
| **Delta from emission logic alone** | **+0.54 ms** | **+638 KB** | **+14,401** |

Two scenarios:

| Prior config | Added cost per scrape |
| --- | --- |
| Already fetching `/pods` (`extra_metadata_labels` or any `*_utilization` metric) | **+0.54 ms** — unmarshal already paid |
| Not fetching `/pods` (the stock default) | **~5.0 ms** — unmarshal + emission |

**In context:** 5 ms per scrape at `collection_interval: 10s` is **0.05% of one core**.
At a 1 s interval it is ~0.5%.

### 6.3 Memory

The ~1.9 MB worst case is **transient garbage**, released immediately after each scrape:

- **No caches, no informers, no retained state were added.**
- The pod list was already allocated and dropped each scrape in configs that fetched `/pods`.
- Expect a few MB of RSS from GC headroom, not a proportional increase.
- Cost is linear in **pods-per-node**, which Kubernetes caps (110 default, 250 max) — bounded.

**Cluster-wide this is a large net win.** It replaces `k8sclusterreceiver`'s informer cache of
*every pod in the cluster* (hundreds of MB at scale) plus a persistent API server watch, with a
bounded per-node poll of localhost.

### 6.4 Output volume

Same 60-pod / 120-container node:

| Configuration | ResourceMetrics | Data points | Δ vs baseline |
| --- | ---: | ---: | ---: |
| Baseline (state metrics off) | 181 | 724 | — |
| State metrics, **without** `status.reason` | 361 | 1,564 | **+116%** |
| State metrics, **all** (current default) | 361 | 2,644 | **+265%** |

Two observations:

1. **`k8s.container.status.reason` dominates** — 1,080 of the 2,644 data points (~40%).
   Disabling it alone drops the increase from +265% to +116%.

2. **ResourceMetrics doubles (181 → 361).** Because state metrics are driven by the pod list
   (§3.3), each pod/container emits a second group carrying a duplicate copy of its resource
   attributes — extra OTLP serialisation and export bytes.

   This was accepted deliberately: merging into the usage groups would silently drop
   `Pending` / `ImagePullBackOff` / `CrashLoopBackOff` pods. A future optimisation could merge
   into the usage group when stats exist and emit standalone groups only for pods that lack them.

### 6.5 Kubelet-side cost

`/pods` serialisation load lands on **the kubelet**, not only the collector. The kubelet serves
from its in-memory pod manager but JSON-marshals full pod objects on every request. Negligible at
10 s+ intervals; measurable on a dense node polled at 1 s.

### 6.6 Tuning

Split into two receiver instances — usage at a fast interval, state at a slow one. State changes
slowly, so 60 s loses almost nothing and cuts both the `/pods` cost and the volume by 6×:

```yaml
receivers:
  kubelet_stats/usage:
    collection_interval: 10s
    metrics:
      k8s.pod.phase: {enabled: false}
      k8s.pod.status_reason: {enabled: false}
      k8s.container.restarts: {enabled: false}
      k8s.container.ready: {enabled: false}
      k8s.container.status.reason: {enabled: false}

  kubelet_stats/state:
    collection_interval: 60s
    # …disable the default usage metrics here
```

---

## 7. Operational Requirements

### 7.1 RBAC — mandatory

Because state metrics are **enabled by default**, `/pods` is fetched on every scrape.
A `/pods` failure returns an error from `scrape()`, which fails the **entire scrape** —
usage metrics included. A collector without this permission goes dark, it does not degrade.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: otel-collector
rules:
  - apiGroups: [""]
    resources: ["nodes/stats"]
    verbs: ["get"]

  # Required by default — pod/container state metrics are enabled by default
  - apiGroups: [""]
    resources: ["nodes/pods"]
    verbs: ["get"]
```

### 7.2 Opting out

```yaml
receivers:
  kubelet_stats:
    metrics:
      k8s.pod.phase: {enabled: false}
      k8s.pod.status_reason: {enabled: false}
      k8s.container.restarts: {enabled: false}
      k8s.container.ready: {enabled: false}
      k8s.container.status.reason: {enabled: false}
```

With all eleven disabled, `needsPodStatus` is false and behaviour reverts exactly to
pre-change: no `/pods` call, no new data points, no extra RBAC.

### 7.3 Enabling the opt-in resource metrics

```yaml
receivers:
  kubelet_stats:
    collection_interval: 10s
    auth_type: serviceAccount
    endpoint: "${env:K8S_NODE_NAME}:10250"
    insecure_skip_verify: true
    metrics:
      k8s.container.cpu_limit: {enabled: true}
      k8s.container.cpu_request: {enabled: true}
      k8s.container.memory_limit: {enabled: true}
      k8s.container.memory_request: {enabled: true}
      k8s.container.storage_limit: {enabled: true}
      k8s.container.storage_request: {enabled: true}
```

---

## 8. Migration Checklist

- [ ] Grant `get` on `nodes/pods` in the DaemonSet ClusterRole **before** rolling out.
- [ ] Confirm pod/container metrics are disabled in `k8sclusterreceiver`. If both emit during a
      rollout window, series **duplicate rather than dedupe** — the resource attribute sets
      differ (§5.2), so backends will double-count.
- [ ] Audit dashboards/alerts for `k8s.node.name`, `container.image.*`, `k8s.pod.qos_class`,
      `k8s.container.status.last_terminated_reason` on these metrics — these attributes are not
      emitted by kubeletstats.
- [ ] Size backend ingest for the volume increase (§6.4) — it now scales per-node rather than
      once cluster-wide.
- [ ] Decide on `k8s.container.status.reason`: keep for k8scluster parity, or disable to cut
      ~40% of the added volume.

---

## 9. Summary

| Dimension | Impact |
| --- | --- |
| **CPU** | ~0.05% of one core at a 10 s interval — negligible |
| **Memory** | A few MB transient; no new caches or informers; bounded by pods-per-node |
| **Cluster-wide memory** | **Large reduction** — removes k8scluster's all-pods informer cache |
| **API server** | **Load removed** — no pod watch |
| **Data point volume** | +265% on pod/container metrics (+116% without `status.reason`) |
| **Metric parity** | Full parity with k8scluster on names, types, units, encodings |
| **Resource attributes** | Narrower than k8scluster (§5.2) — intentional |
| **Breaking change** | Yes — `nodes/pods` RBAC becomes mandatory |
