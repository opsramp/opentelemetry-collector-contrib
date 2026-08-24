# Impact Analysis: opentelemetry-collector-contrib v0.152.0 → v0.159.0

**Fork:** `opsramp/opentelemetry-collector-contrib`
**Branch:** `release/v0.159.0` (merged from `upstream/release/v0.159.x`)
**Core collector version:** `v1.65.0` (stable) / `v0.159.0` (unstable)
**Analysis scope:** packages consumed by `agent/opsramp/go.mod`

---

## 1. Removed Packages / Import Paths (CRITICAL)

| Old Import Path | Replacement | Agent Uses? | Action Required |
|---|---|---|---|
| `receiver/jmxreceiver` | **Deleted upstream** | **Yes** | None — retained in our fork, deps bumped to v0.159.0. Fork now permanently owns this component. |
| `extension/observer/kafkatopicsobserver` | Deleted upstream | No | Removed from `builder-config.yaml` |
| `go.opentelemetry.io/collector/confmap/xconfmap` | `go.opentelemetry.io/collector/confmap` | No (test-only) | `xconfmap.Validate` → `confmap.Validate` (already fixed in contrib) |

All 63 contrib import paths referenced from agent Go source were verified to still exist at v0.159.0.

---

## 2. Component Renames

| Old Name | New Name | Alias Maintained? | Agent Uses? |
|---|---|---|---|
| `apachespark` | `apache_spark` | Yes | Yes |
| `prometheusremotewrite` | `prometheus_remote_write` | Yes | Yes |
| `resourcedetection` | `resource_detection` | Yes | No |
| `sqlquery` | `sql_query` | Yes | No |
| `deltatocumulative` | `delta_to_cumulative` | Yes | No |
| `deltatorate` | `delta_to_rate` | Yes | No |
| `cumulativetodelta` | `cumulative_to_delta` | Yes | No |
| `loadbalancing` | `load_balancing` | Yes | No |
| `otlpjsonfile` | `otlp_json_file` | Yes | No |
| `webhookevent` | `webhook_event` | Yes | No |
| `metricsaslogs` | `metrics_as_logs` | Yes | No |
| `awscloudwatch` | `aws_cloudwatch` | Yes | No |
| `envoyals` | `envoy_als` | Yes | No |
| `spanpruning` | `span_pruning` | Yes | No |
| `awsecsattributes` | `aws_ecs_attributes` | **NO** (in development) | No |

No agent-used component lost its deprecated alias.

### Component type names verified (agent-used)

| Package | Type | Deprecated Alias |
|---|---|---|
| `receiver/hostmetricsreceiver` | `host_metrics` | `hostmetrics` |
| `receiver/filelogreceiver` | `file_log` | `filelog` |
| `receiver/k8sclusterreceiver` | `k8s_cluster` | — (renamed before v0.152) |
| `receiver/k8seventsreceiver` | `k8s_events` | — (renamed before v0.152) |
| `receiver/k8sobjectsreceiver` | `k8s_objects` | `k8sobjects` |
| `receiver/kubeletstatsreceiver` | `kubelet_stats` | `kubeletstats` |
| `receiver/windowseventlogreceiver` | `windows_event_log` | `windowseventlog` |
| `receiver/kafkametricsreceiver` | `kafka_metrics` | `kafkametrics` |
| `receiver/apachesparkreceiver` | `apache_spark` | `apachespark` |
| `processor/k8sattributesprocessor` | `k8s_attributes` | `k8sattributes` |
| `exporter/prometheusremotewriteexporter` | `prometheus_remote_write` | `prometheusremotewrite` |

---

## 3. Removed Feature Gates

| Component | Gate ID | Was Behavior | Impact |
|---|---|---|---|
| `k8s_attributes` | `k8sattr.labelsAnnotationsSingular.allow` | Stable / on | Permanently on — no action if agent never disabled it |
| `pkg/fileconsumer` | `filelog.decompressFingerprint` | Stable / on | Permanently on |
| `pkg/fileconsumer` | `filelog.protobufCheckpointEncoding` | Alpha → **Beta, on by default** | Checkpoint encoding changed; **downgrade path is not clean** |

---

## 4. Breaking Changes Affecting OpsRamp Agent

| Component | Change | Risk | Action Required |
|---|---|---|---|
| `host_metrics` | `cpu` attribute now **opt-in** on `system.cpu.time` / `system.cpu.utilization`; metrics aggregated across logical CPUs (#49161) | **High** | Agent's `k8s_metrics_validator.py` asserts `system.cpu.utilization`. Re-enable the attribute (see snippet below) to restore per-core output |
| `pkg/ottl` | Datapoint context setters now **error** on incompatible datapoint types instead of silently no-opping (#48384) | Medium | Audit OTTL statements in transform/filter configs |
| `pkg/ottl` | `is_root_span` API removed (#50093) | Low | Not used by agent (verified) |
| `pkg/fileconsumer` | Checkpoint encoding switched to protobuf by default (#49387) | Medium | Rolling back to v0.152 will invalidate checkpoints → possible log re-read |
| `haproxy` | Metric units changed to comply with UCUM spec (#49453) | Low | Dashboards/alerts on haproxy units may need updates |

### Restoring per-core CPU metrics

```yaml
receivers:
  hostmetrics:
    scrapers:
      cpu:
        metrics:
          system.cpu.time:
            attributes: [cpu, state]
          system.cpu.utilization:
            attributes: [cpu, state]
```

---

## 5. Deprecated APIs in Use

| Package | Deprecated API | Replacement | Deadline |
|---|---|---|---|
| `file_log` | Implicit `ordering_criteria.top_n = 1` when `sort_by` is set | Set `top_n` explicitly; `filelog.requireExplicitTopN` gate enforces it | Gate expected to become default in a future release |
| `routing` connector | `request` context / `request["key"]` syntax | `otelcol.client.metadata["key"][0]` (HTTP) / `otelcol.grpc.metadata["key"][0]` (gRPC) | Not announced |
| `k8s_attributes` | `deployment_name_from_replicaset` | ReplicaSet-name heuristic (now the default) | Not announced |
| `resource_detection` | `k8snode` detector | `k8s_api` detector (rename config key too) | Not announced |

---

## 6. Disabled-by-Default Changes

| Component | What Changed | Previous Default | New Default | Agent Impact |
|---|---|---|---|---|
| `host_metrics` | `cpu` attribute on CPU time/utilization | enabled | **disabled** | Per-core CPU breakdown lost unless re-enabled |
| `signalfx` exporter | per-core `cpu.*` translations | enabled | disabled | None (not used by agent) |
| `signalfx` exporter | `cpu.utilization_per_core` translation | enabled | disabled | None (not used by agent) |

---

## 7. Required Agent go.mod Updates

- All 25 `github.com/open-telemetry/opentelemetry-collector-contrib/*` deps: `v0.152.0` → `v0.159.0`
- All 10 `go.opentelemetry.io/collector/*` stable deps: `v1.58.0` → `v1.65.0`
- Unstable `go.opentelemetry.io/collector/*` deps: `v0.152.0` → `v0.159.0`

> **Caution:** `go.opentelemetry.io/collector/config/configgrpc` graduated from the `v0.x` line to **`v1.65.0`**.
> A blind `v0.152.0` → `v0.159.0` search-and-replace will fail to resolve
> (`unknown revision config/configgrpc/v0.159.0`). Verify each collector module against
> a module already updated by upstream before bumping.

**No import path changes are required in agent source code.**

---

## 8. New Features Available

- `k8s_objects`: top-level `interval` field as a fallback default pull interval for all pull-mode objects (#48452)
- `k8s_events` (upstream): `dedup_interval` MODIFIED-event throttling — **not adopted** (OpsRamp receiver logic kept)
- `k8s_attributes` (upstream): configurable `pod_delete_grace_period` — **not adopted** (hardcoded to the previous 120s default)
- `pkg/stanza`: EVT_HANDLE leak fix in the Windows Event Log receiver when no checkpoint has been persisted (#47194)
- `processor/schema`: support for Schema v2 file formats (`manifest/2.0`, `resolved/2.0`)

---

## 9. Risk Assessment Summary

**Overall risk: Medium-High**

Recommended testing focus areas, in priority order:

1. **hostmetrics CPU metrics** — highest risk; validate `system.cpu.*` output shape end-to-end against dashboards and `k8s_metrics_validator.py`
2. **filelog checkpoint migration** — verify no log duplication or loss on collector restart after upgrade
3. **k8seventsreceiver / k8sobjectsreceiver** — OpsRamp logic retained, but `k8sinventory` `Observer.Start` signature changed to `(chan struct{}, error)`
4. **k8sattributesprocessor** — pod delete grace period path; note `watch_sync_period` remains 30m (ours) rather than upstream's new 5m default
5. **OTTL statements** in transform/filter processors — previously silent no-ops now surface as errors

---

## 10. Merge Conflict Resolution

The merge of `upstream/release/v0.159.x` produced **27 conflicted files**.

### 10.1 Summary by resolution strategy

| Resolution | Files | Rationale |
|---|---|---|
| **Kept OURS** (OpsRamp logic preserved) | 22 | 14 protected-package files + 8 jmxreceiver files |
| **Took THEIRS** (upstream) | 4 | Dependency manifests only — no logic |
| **Manual merge** | 1 | `builder-config.yaml` — upstream base + OpsRamp components re-added |

### 10.2 Full conflicted-file list

#### Kept OURS — protected packages (14 files)

| # | File | Package | What ours preserves |
|---|---|---|---|
| 1 | `receiver/k8seventsreceiver/config.go` | k8seventsreceiver | `EventType` constants, `ExcludeNamespaces`, `IncludeInvolvedObject`; rejects upstream `DedupInterval` + `APIConfig` field rename |
| 2 | `receiver/k8seventsreceiver/receiver.go` | k8seventsreceiver | `excludedNSSet` / `excludedReasonsSet` filtering, `toSet()`, `computeNamespaceFilter()`, attribute-returning `allowEvent`; rejects upstream dedup cache + telemetry builder |
| 3 | `receiver/k8seventsreceiver/receiver_test.go` | k8seventsreceiver | `TestComputeNamespaceFilter` and namespace-exclusion tests |
| 4 | `receiver/k8seventsreceiver/testdata/config.yaml` | k8seventsreceiver | OpsRamp config shape |
| 5 | `receiver/k8sobjectsreceiver/receiver.go` | k8sobjectsreceiver | Custom `startPull` / `startWatch` with pagination (`PageLimit`, `PageInterval`), namespace discovery, `k8sinventory` mode handling |
| 6 | `processor/k8sattributesprocessor/config.go` | k8sattributesprocessor | Rejects upstream `PodDeleteGracePeriod` config field |
| 7 | `processor/k8sattributesprocessor/factory.go` | k8sattributesprocessor | `WatchSyncPeriod: 30m` (upstream lowered to 5m) |
| 8 | `processor/k8sattributesprocessor/options.go` | k8sattributesprocessor | OpsRamp option set |
| 9 | `processor/k8sattributesprocessor/processor.go` | k8sattributesprocessor | `semconv/v1.5.0` + `conventions/v1.40.0`, `processTraceResources`, Redis client path, single-arg `processResource` |
| 10 | `pkg/stanza/fileconsumer/file.go` | pkg/stanza | Custom `makeReader` (fingerprint dedup) that upstream removed |
| 11 | `pkg/stanza/fileconsumer/matcher/matcher.go` | pkg/stanza | `RefreshInterval`, `defaultOrderingCriteriaTopN`, non-pointer `TopN int` |
| 12 | `pkg/stanza/fileconsumer/matcher/matcher_test.go` | pkg/stanza | Matches our matcher API |
| 13 | `pkg/stanza/fileconsumer/matcher/internal/finder/finder_test.go` | pkg/stanza | Rejects upstream case-sensitivity featuregate test |
| 14 | `pkg/stanza/operator/input/windows/input.go` | pkg/stanza | OpsRamp raw/simple event processing path and error logging |

#### Kept OURS — jmxreceiver (8 files, `modify/delete` conflicts)

Upstream **deleted the entire `receiver/jmxreceiver`**. The whole directory (19 files) was
restored from our branch since the agent depends on it.

| # | File | Conflict type |
|---|---|---|
| 15 | `receiver/jmxreceiver/config.go` | deleted upstream, modified in HEAD |
| 16 | `receiver/jmxreceiver/config_test.go` | deleted upstream, modified in HEAD |
| 17 | `receiver/jmxreceiver/factory.go` | deleted upstream, modified in HEAD |
| 18 | `receiver/jmxreceiver/factory_test.go` | deleted upstream, modified in HEAD |
| 19 | `receiver/jmxreceiver/go.mod` | deleted upstream, modified in HEAD |
| 20 | `receiver/jmxreceiver/receiver.go` | deleted upstream, modified in HEAD |
| 21 | `receiver/jmxreceiver/receiver_test.go` | deleted upstream, modified in HEAD |
| 22 | `receiver/jmxreceiver/supported_jars.go` | deleted upstream, modified in HEAD |

> The remaining ~11 jmxreceiver files (Makefile, README.md, metadata.yaml, testdata,
> `internal/subprocess/`, generated files) were staged as deletions by git and restored
> via `git checkout HEAD -- receiver/jmxreceiver/`.

#### Took THEIRS — dependency manifests only (4 files)

| # | File | Rationale |
|---|---|---|
| 23 | `pkg/stanza/go.sum` | Hashes only; `replace` directives identical on both sides |
| 24 | `processor/k8sattributesprocessor/go.mod` | Version bumps; all 8 `replace` directives verified byte-identical to ours before taking theirs. OpsRamp-only deps (`redis/go-redis/v9`, `hashicorp/golang-lru/v2`, `go-pkgz/expirable-cache/v3`) re-added automatically by `go mod tidy` |
| 25 | `processor/k8sattributesprocessor/go.sum` | Hashes only |
| 26 | `receiver/hostmetricsreceiver/go.mod` | Version bumps only (`v1.58.0`→`v1.65.0`, `v0.152.0`→`v0.159.0`) |

#### Manual merge (1 file)

| # | File | Resolution |
|---|---|---|
| 27 | `cmd/otelcontribcol/builder-config.yaml` | Took upstream as base (all v0.159.0 versions + new upstream components), then re-added the 5 OpsRamp components at v0.159.0 plus `jmxreceiver` and `signalfxreceiver` |

Re-added entries:
- `exporter/opsrampdebugexporter v0.159.0`
- `exporter/opsrampotlpexporter v0.159.0`
- `processor/opsrampmetricsfilterprocessor v0.159.0`
- `processor/opsrampk8sobjectsprocessor v0.159.0`
- `processor/scrubbingprocessor v0.159.0`
- `receiver/jmxreceiver v0.159.0` (retained fork component)
- `receiver/signalfxreceiver v0.159.0` (dropped upstream, directory still exists)

Dropped: `extension/observer/kafkatopicsobserver` (directory removed upstream).

### 10.3 Post-merge reverts (auto-merged files reverted to ours)

These files merged cleanly but pulled in upstream code that referenced APIs our kept
source does not have. They were restored from `release/v0.152.0`:

| File | Reason |
|---|---|
| `receiver/k8seventsreceiver/config_test.go` | Referenced upstream `DedupInterval` |
| `receiver/k8seventsreceiver/README.md` | Documented dedup feature |
| `receiver/k8seventsreceiver/config.schema.yaml` | Declared `dedup_interval` |
| `receiver/k8seventsreceiver/metadata.yaml` | Declared dedup telemetry metrics |
| `processor/k8sattributesprocessor/config_test.go` | Referenced upstream `PodDeleteGracePeriod` |
| `processor/k8sattributesprocessor/testdata/config.yaml` | Contained `pod_delete_grace_period` |
| `processor/k8sattributesprocessor/README.md` | Documented `pod_delete_grace_period` |
| `processor/k8sattributesprocessor/config.schema.yaml` | Declared `pod_delete_grace_period` |
| `pkg/stanza/fileconsumer/config_test.go` | Expected pointer `TopN *int` |
| `pkg/stanza/fileconsumer/file_test.go` | Expected pointer `TopN *int` |

### 10.4 Files deleted (upstream dedup-only additions)

Removed to stay consistent with the retained OpsRamp k8sevents receiver:

- `receiver/k8seventsreceiver/e2e_test.go`
- `receiver/k8seventsreceiver/documentation.md`
- `receiver/k8seventsreceiver/testdata/e2e/` (entire tree)
- `receiver/k8seventsreceiver/internal/metadata/generated_telemetry.go`
- `receiver/k8seventsreceiver/internal/metadata/generated_telemetry_test.go`
- `receiver/k8seventsreceiver/internal/metadatatest/` (entire directory)

### Build fixes applied

| File | Fix |
|---|---|
| `receiver/k8seventsreceiver/receiver.go` | `observer.Start` now returns `(chan struct{}, error)` |
| `processor/k8sattributesprocessor/processor.go` | New `podDeleteGracePeriod` param on `kube.New`; passed the previous `120s` default |
| `receiver/k8seventsreceiver/config_test.go` | `xconfmap.Validate` → `confmap.Validate` |
| `processor/k8sattributesprocessor/config_test.go` | `xconfmap.Validate` → `confmap.Validate` |
| `receiver/jmxreceiver/config_test.go` | `xconfmap.Validate` → `confmap.Validate` |

Upstream's k8sevents dedup support files (e2e tests, generated telemetry, `documentation.md`,
`metadatatest/`) were removed to stay consistent with the retained OpsRamp receiver implementation.

### Build verification (all passing `go build ./...` and `go vet ./...`)

`internal/k8sinventory`, `receiver/k8seventsreceiver`, `receiver/k8sobjectsreceiver`,
`pkg/stanza`, `processor/k8sattributesprocessor`, `receiver/jmxreceiver`,
`receiver/prometheusreceiver`, `receiver/hostmetricsreceiver`,
`processor/opsrampk8sobjectsprocessor`, `processor/opsrampmetricsfilterprocessor`,
`processor/scrubbingprocessor`, `exporter/opsrampdebugexporter`, `exporter/opsrampotlpexporter`

---

## 11. Post-Merge Behavioural Audit (silent breaks found & fixed)

`go build` on macOS is **not sufficient** to validate this merge. Two defects compiled
cleanly but were behaviourally broken. Both are fixed.

### 11.1 Windows Event Log was completely broken (CRITICAL — fixed)

`input.go` was kept as ours, but its collaborators (`subscription.go`, `publisher.go`,
`publishercache.go`, `api.go`, `config_windows.go`) auto-merged to upstream's refactored
APIs from the EVT_HANDLE leak fix (#47364), which moved bookmark/retry ownership from
`Subscription` into `Input` and changed `Subscription.Read` from 3 to 2 return values.

Because those files are `//go:build windows`, a macOS `go build ./...` never compiled them.
`GOOS=windows go build` produced 7 errors:

```
input.go:273: assignment mismatch: 3 variables but i.subscription.Read returns 2 values
input.go:156: not enough arguments in call to subscription.Open
config_windows.go:75: unknown field eventDrivenScraping in struct literal of type Input
...
```

**Resolution:** restored `pkg/stanza/operator/input/windows/` to ours, then **ported the
EVT_HANDLE leak fix (#47364) forward** onto our customised `input.go`:
- Took `subscription.go` + `subscription_test.go` at commit `94e2e402772` (verified our
  v0.152.0 `subscription.go` was byte-identical to that commit's parent, so the port is exact)
- Added `Input.readWithRetry` — moves RPC_S_INVALID_BOUND retry from `Subscription` to `Input`
- `readBatch` now calls `readWithRetry` instead of the 3-value `subscription.Read`
- Updated `input_test.go` to put reopen state (`startAt`, `channel`) on `Input` rather than `Subscription`

Upstream's EVTX file support (#48047) and event-driven scraping (#48463) were **not** adopted.

### 11.2 k8sobjectsreceiver silently ignored `initial_delay` (fixed)

`config.go` auto-merged and accepted + validated upstream's new `initial_delay`, but our
kept `receiver.go` had no code reading it — a user setting `initial_delay: 10m` would get
no delay and no error.

**Resolution:** implemented it in our receiver rather than dropping the config field:
- `InitialDelay` field + validation in `config.go` (negative, watch-mode, and `>= interval` rejected)
- Delay honoured in `startPull`, interruptible by both `stopperChan` and `ctx.Done()`
- Added `TestPullObjectInitialDelay` and `TestPullObjectShutdownDuringInitialDelay`

> While implementing this, a related latent bug was found: `K8sObjectsConfig.DeepCopy()`
> silently drops `PageLimit` and `PageInterval`, so those config options are accepted but
> never reach `startPull` (it falls back to a hardcoded 500). **Not fixed here** — changing
> it would alter current pagination behaviour. Recommend tracking separately.

### 11.3 Regression testing vs the v0.152.0 baseline

Test suites were run side-by-side against a git worktree of the pre-upgrade branch:

| Package | HEAD failures | Baseline failures | Verdict |
|---|---|---|---|
| `receiver/k8seventsreceiver` | `TestK8sEventToLogData`, `TestLoadConfig` | identical | No regression (pre-existing) |
| `processor/k8sattributesprocessor` | 5 tests | identical | No regression (pre-existing) |
| `processor/k8sattributesprocessor/internal/kube` | **pass** | pass | Confirms the `podDeleteGracePeriod` fix |
| `receiver/k8sobjectsreceiver` | 4 tests (after fix) | identical 4 | No regression (pre-existing) |

`pkg/stanza/fileconsumer` failures on HEAD (`TestRestartOffsets`, `TestHeaderPersistance`,
`TestMultiFileSort*`, `TestFileMovedWhileOff_BigFiles`, plus `internal/checkpoint` encoding
tests) have **not yet been baselined** — the comparison run was cancelled. This is the one
outstanding verification gap.

### 11.4 Cross-platform build verification

All 13 packages built for `GOOS=windows`, `GOOS=linux`, `GOOS=darwin` — **0 failures**.
This sweep should be part of the standard upgrade procedure; a native-only build hides
platform-gated breakage.

### 11.5 Verified-benign upstream drift

Files our kept sources depend on that silently took upstream changes, reviewed individually:

| File | Change | Verdict |
|---|---|---|
| `fileconsumer/config.go` | New opt-in `file_cache_advise` field (default false) | Benign |
| `fileconsumer/matcher/internal/filter/sort.go` | Precompiled strptime parser | Perf only, same semantics |
| `fileconsumer/file.go` (upstream changes rejected) | gofumpt + function reordering only | No behaviour rejected |
| `fileconsumer/matcher/matcher.go` (rejected) | `requireExplicitTopN` gate, off by default | Our behaviour == upstream default |

---

## 12. OpsRamp Customisation Preservation Audit

**Question answered:** did any OpsRamp custom change get removed or overridden by upstream?

**Method (not opinion — measured):** the set of OpsRamp-customised files is
`diff(upstream v0.152.0 tag → our release/v0.152.0)` = **228 files** (excluding `go.mod`/`go.sum`).
Each was then checked at HEAD.

### 12.1 Headline result

| Outcome | Count |
|---|---|
| Preserved **byte-identical** to v0.152.0 | **203** |
| Changed (all reviewed individually — see 12.3) | 25 |
| **Deleted / lost** | **0** |

For the 16 changed non-`go.mod` files, the *OpsRamp delta itself* was re-extracted at HEAD
(`diff(upstream v0.159 → HEAD)`) and compared against the original delta
(`diff(upstream v0.152 → our v0.152)`):

| Delta comparison | Count |
|---|---|
| Delta **identical** — customisation untouched | 9 |
| Delta differs — **only** because of a deliberate edit made during this upgrade | 7 |

**No OpsRamp customisation was lost, reverted, or overwritten.**

### 12.2 Preserved byte-identical, by area

| Area | Custom files preserved | Changed |
|---|---|---|
| `exporter/opsrampdebugexporter` | 56 | 1 (`go.mod` only) |
| `receiver/hostmetricsreceiver` | 30 | 3 |
| `pkg/stanza` | 23 | 6 |
| `exporter/opsrampotlpexporter` | 15 | 1 (`go.mod` only) |
| `processor/opsrampk8sobjectsprocessor` | 13 | 1 (`go.mod` only) |
| `processor/scrubbingprocessor` | 11 | 1 (`go.mod` only) |
| `processor/opsrampmetricsfilterprocessor` | 11 | 1 (`go.mod` only) |
| `processor/k8sattributesprocessor` | 11 | 3 |
| `extension/opsramplogsdbstorage` | 10 | **0** |
| `receiver/jmxreceiver` | 6 | 2 |
| `receiver/k8seventsreceiver` | 5 | 2 |
| `cmd/otelcontribcol` | 2 | 1 |
| `receiver/k8sobjectsreceiver` | 1 | 2 |
| `internal/k8sinventory` | 1 (`config.go`) | 0 |
| `receiver/prometheusreceiver` | 1 | 0 |
| `.github/*`, `deploy.yaml`, `go.mod` | 5 | 1 (`Makefile`) |

> **All 5 OpsRamp-only components plus `opsramplogsdbstorage` have 100% of their source
> byte-identical to v0.152.0.** Their only change is the `go.mod` version bump.

### 12.3 The 25 changed files — every one accounted for

**9 dependency manifests** (`go.mod`) — version bumps only, required by the upgrade:
`opsrampdebugexporter`, `opsrampotlpexporter`, `opsrampk8sobjectsprocessor`,
`opsrampmetricsfilterprocessor`, `scrubbingprocessor`, `k8sattributesprocessor`,
`pkg/stanza`, `hostmetricsreceiver`, `jmxreceiver`.

**9 files where the OpsRamp delta is byte-identical** (upstream changed surrounding code
only; our customisation is intact):

| File | OpsRamp customisation confirmed still present |
|---|---|
| `processor/k8sattributesprocessor/internal/kube/kube.go` | `AddOnMetadata` struct |
| `receiver/hostmetricsreceiver/factory.go` | `groupprocessscraper.NewFactory()` registration |
| `receiver/hostmetricsreceiver/README.md` | groupprocess docs |
| `pkg/stanza/fileconsumer/config.go` | `InitialBufferSize` default fallback |
| `pkg/stanza/fileconsumer/internal/archive/archive.go` | `writeArchive` without batched storage ops |
| `pkg/stanza/fileconsumer/matcher/internal/filter/sort.go` | `SortMtime(ascending, duration)` + `maxTime`/`hasLimit` age cutoff |
| `pkg/stanza/fileconsumer/matcher/internal/filter/sort_test.go` | matching tests |
| `Makefile` | OpsRamp targets |
| *(`cmd/otelcontribcol/builder-config.yaml` — see below)* | all OpsRamp components carried over |

**7 files changed by a deliberate edit during this upgrade** — all additive, no OpsRamp
logic removed:

| File | Edit | OpsRamp logic touched? |
|---|---|---|
| `receiver/k8seventsreceiver/receiver.go` | `observer.Start` 2-value return + gofmt alignment | No — `excludedNSSet`, `excludedReasonsSet`, `computeNamespaceFilter`, `allowEvent` all intact |
| `processor/k8sattributesprocessor/processor.go` | Added `defaultPodDeleteGracePeriod = 120s` const, passed to `kube.New` | No — preserves the exact previous hardcoded value |
| `receiver/k8sobjectsreceiver/config.go` | Added `InitialDelay` field + validation + `DeepCopy` entry | No — purely additive |
| `receiver/k8sobjectsreceiver/receiver.go` | Added interruptible delay at top of `startPull` | No — purely additive |
| `pkg/stanza/operator/input/windows/input.go` | Added `readWithRetry`, switched `readBatch` to it (EVT_HANDLE leak fix port) | No — replaces the old 3-value `subscription.Read` call only |
| `receiver/k8seventsreceiver/config_test.go` | `xconfmap.Validate` → `confmap.Validate` | No — API rename only |
| `receiver/jmxreceiver/config_test.go` | `xconfmap.Validate` → `confmap.Validate` | No — API rename only |

### 12.4 builder-config.yaml

Version-normalised comparison of all OpsRamp/fork-specific entries between v0.152.0 and HEAD:
**no difference** — every OpsRamp component, plus the fork-retained `jmxreceiver` and
`signalfxreceiver`, carried over.

Note: `opsramplogsdbstorage` is **not** referenced in `builder-config.yaml` at HEAD — but it
was not referenced at v0.152.0 either, so this is a pre-existing state, not a regression.
Flagged in case it is expected to be wired in.

---

## 13. Pre-existing Issues (NOT caused by this upgrade)

Both confirmed present on `release/v0.152.0`:

- `pkg/stanza/operator/input/windows/xml_test.go:79` — `[]string` vs `string` type mismatch on
  `EventXML.Keywords` (12 occurrences). **This breaks compilation of the entire windows test
  package**, meaning the Windows Event Log unit tests have not been running at all. Left
  unfixed (out of scope), but it should be prioritised — it is why 11.1 went undetected.
- `exporter/opsrampotlpexporter/otlp_test.go:738` — `getAuthToken` arity mismatch
  (`SecuritySettings` vs `SecuritySettings, TLSOptions`)

Additional note: `receiver/signalfxreceiver` was dropped from upstream's `builder-config.yaml`
but the directory still exists; the entry was retained in our config.
