---
name: k8s-events-receiver
description: Use when editing, debugging, or enhancing the k8s_events receiver in opentelemetry-collector-contrib. MANDATORY for all changes to receiver/k8seventsreceiver/, including config, factory, receiver lifecycle, event-to-log conversion, event filtering (allowEvent), namespace watching (including exclude_namespaces filter), leader election integration, or test updates. FORBIDDEN patterns include blocking in watch handler callbacks, missing ObsReport instrumentation, bypassing allowEvent timestamp guard, modifying generated files manually, breaking the factory→newReceiver→Start lifecycle, using typed client-go informers (replaced by dynamic client + watch observer), duplicating set-building logic (use toSet() helper). REQUIRED patterns are dynamicfake.NewSimpleDynamicClient for test mocking, ObsReport Start/End bracketing, event timestamp priority (EventTime > LastTimestamp > FirstTimestamp), consumertest.LogsSink for test assertions, goleak.VerifyTestMain for goroutine leak detection, k8sinventory/watch.Observer for namespace watching, sync.Mutex protection on stopperChanList, sync.WaitGroup for goroutine lifecycle, toSet() for building string sets, excludedNSSet built once in newReceiver() for O(1) allowEvent() lookups.
dependsOn: []
---

# K8s Events Receiver — Development Skill

## THE MANDATE

The `k8s_events` receiver in `opentelemetry-collector-contrib` watches Kubernetes API server events via a **dynamic client** and the `k8sinventory/watch.Observer` abstraction, converts them to OpenTelemetry Log records, and supports **optional leader election** via the `k8sleaderelector` extension. Every change MUST respect the OTEL collector component lifecycle (`Start`/`Shutdown`), the event filtering pipeline, the plog data model mapping, and the opentelemetry-collector-contrib testing conventions.

> **Source location:** `receiver/k8seventsreceiver/` in opentelemetry-collector-contrib
> **Component type:** `k8s_events` (Logs receiver, alpha stability)
> **Codeowners:** @dmitryax, @TylerHelmuth, @ChrsMark
> **Distributions:** contrib, k8s

---

## ARCHITECTURE OVERVIEW

```
┌─────────────────────────────────────────────────────────────────────┐
│                    k8s_events Receiver                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  factory.go                                                         │
│  ├── NewFactory() → registers "k8s_events" type                    │
│  ├── createDefaultConfig() → Config{AuthType: serviceAccount}      │
│  └── createLogsReceiver() → newReceiver(settings, config, consumer)│
│                                                                     │
│  receiver.go                                                        │
│  ├── k8seventsReceiver struct                                      │
│  │   ├── config          *Config                                   │
│  │   ├── excludedNSSet   map[string]struct{} (built in newReceiver, nil when no exclusions) │
│  │   ├── settings        receiver.Settings                         │
│  │   ├── logsConsumer    consumer.Logs                             │
│  │   ├── stopperChanList []chan struct{}                            │
│  │   ├── startTime       time.Time                                 │
│  │   ├── ctx/cancel      context.Context/CancelFunc               │
│  │   ├── obsrecv         *receiverhelper.ObsReport                 │
│  │   ├── mu              sync.Mutex  (protects stopperChanList)    │
│  │   ├── client          dynamic.Interface                         │
│  │   └── wg              sync.WaitGroup (goroutine lifecycle)      │
│  │                                                                  │
│  ├── Start(ctx, host)                                              │
│  │   ├── Creates context with cancel                               │
│  │   ├── Gets dynamic client via config.getDynamicClient()         │
│  │   ├── If K8sLeaderElector configured:                           │
│  │   │   ├── Resolves extension from host.GetExtensions()          │
│  │   │   ├── Registers OnLeading → startWatchers()                 │
│  │   │   └── Registers OnStopping → Shutdown()                     │
│  │   └── Else: starts immediately → startWatchers()                │
│  │                                                                  │
│  ├── Shutdown(ctx)                                                  │
│  │   ├── Cancels context (if cancel != nil)                        │
│  │   ├── Mutex-protected: closes all stopper channels, nils list   │
│  │   └── wg.Wait() — waits for all goroutines to finish            │
│  │                                                                  │
│  ├── startWatchers()                                               │
│  │   ├── Defines eventsGVR (v1/events)                             │
│  │   ├── Creates watchobserver.Observer with:                      │
│  │   │   ├── IncludeInitialState: false                            │
│  │   │   └── Exclude: {Deleted: true}                              │
│  │   ├── Observer callback: Unstructured → corev1.Event conversion │
│  │   │   └── Calls handleEvent(ev)                                 │
│  │   └── observer.Start(ctx, &wg) → returns stopper channel        │
│  │                                                                  │
│  ├── handleEvent(ev)                                               │
│  │   ├── allowEvent(ev) → filter check                             │
│  │   ├── k8sEventToLogData(logger, ev, version, attributes)        │
│  │   └── obsrecv.Start/End + consumer.ConsumeLogs()                │
│  │                                                                  │
│  ├── allowEvent(ev) → (attributes []KeyValue, allow bool)         │
│  │   ├── Rejects events older than receiver startTime              │
│  │   ├── Filters by EventTypes (Normal/Warning)                    │
│  │   ├── ExcludeNamespaces gate (O(1) via excludedNSSet map)       │
│  │   ├── Filters by IncludeInvolvedObject kind + reasons           │
│  │   └── Returns custom attributes from matched reason             │
│  │                                                                  │
│  ├── getEventTimestamp(ev) → time.Time  (local wrapper)            │
│  │   └── Delegates to k8sinventory.GetEventTimestamp()             │
│  │   └── Priority: EventTime > LastTimestamp > FirstTimestamp       │
│  ├── computeNamespaceFilter(namespaces, excludeNamespaces) []string │
│  │   └── Set difference: removes excludeNamespaces from namespaces  │
│  │   └── Used in startWatchers() Rule 4 (both lists non-empty)     │
│  │   └── Returns nil when all namespaces excluded (caller guards)  │
│  └── toSet(items []string) map[string]struct{}                      │
│      └── Shared helper — used by newReceiver + computeNamespaceFilter│
│                                                                     │
│  k8s_event_to_logdata.go                                           │
│  └── k8sEventToLogData(logger, ev, version, attributes) → plog.Logs│
│      ├── Scope: name=metadata.ScopeName, version=buildInfo.Version │
│      ├── Resource attributes (node, object kind/name/uid/ns)       │
│      ├── Dynamic "k8s.<kind>.name" + "resourceName" attributes     │
│      ├── Log body = ev.Message                                     │
│      ├── Severity: "normal"→Info, "warning"→Warn,                  │
│      │             "error"→Error, "critical"→Fatal                  │
│      └── Log attributes (type, reason, action, count,              │
│          k8s.namespace.name, etc.)                                  │
│                                                                     │
│  config.go                                                          │
│  ├── Config struct (embeds k8sconfig.APIConfig)                    │
│  │   ├── Namespaces         []string                               │
│  │   ├── ExcludeNamespaces  []string  (namespace exclusion list)   │
│  │   ├── EventTypes         []EventType                            │
│  │   ├── IncludeInvolvedObject map[string]InvolvedObjectProperties │
│  │   ├── K8sLeaderElector *component.ID                            │
│  │   ├── makeClient func (for test mocking — typed client)         │
│  │   └── makeDynamicClient func (for test mocking — dynamic client)│
│  ├── InvolvedObjectProperties → IncludeReasons []ReasonProperties  │
│  ├── ReasonProperties → Name string + Attributes []KeyValue        │
│  ├── KeyValue → Key string + Value any                             │
│  ├── Validate() → checks EventTypes + APIConfig.Validate()        │
│  ├── getK8sClient() → typed k8s.Interface                         │
│  └── getDynamicClient() → dynamic.Interface                        │
│                                                                     │
│  metadata.yaml                                                      │
│  └── type: k8s_events, stability: alpha [logs]                     │
│                                                                     │
│  internal/metadata/generated_status.go  (DO NOT EDIT)              │
│  └── Type, ScopeName, LogsStability constants                      │
│                                                                     │
│  EXTERNAL DEPENDENCIES:                                             │
│  ├── internal/k8sinventory/watch/observer.go                       │
│  │   └── Watch observer: manages K8s API watches per namespace     │
│  ├── internal/k8sinventory/utils.go                                │
│  │   └── GetEventTimestamp() — shared timestamp priority logic     │
│  └── extension/k8sleaderelector/                                   │
│      └── LeaderElection interface for HA deployments               │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## FILE MAP

| File | Purpose | Edit Frequency |
|------|---------|----------------|
| `receiver.go` | Receiver struct, Start/Shutdown, watch lifecycle, event filtering (`allowEvent`), `handleEvent` | High — bugs, new filters, lifecycle changes |
| `config.go` | Config struct, validation, K8s client factory, filter types (`EventType`, `InvolvedObjectProperties`, `ReasonProperties`, `KeyValue`) | High — new config fields, validation rules |
| `factory.go` | Factory registration, default config, `createLogsReceiver` | Low — rarely changes |
| `k8s_event_to_logdata.go` | Event-to-plog conversion, resource/log attributes, severity mapping | Medium — new attributes, mapping changes |
| `doc.go` | Package doc + `//go:generate mdatagen metadata.yaml` | Rare |
| `metadata.yaml` | Component metadata (type, stability, codeowners, test config) | Rare — stability promotions |
| `internal/metadata/generated_status.go` | **GENERATED — DO NOT EDIT** | Never manual |
| `generated_component_test.go` | **GENERATED — DO NOT EDIT** | Never manual |
| `generated_package_test.go` | **GENERATED — goroutine leak test via goleak** | Never manual |
| `receiver_test.go` | Tests for receiver lifecycle, handleEvent, allowEvent, timestamp | High |
| `k8s_event_to_logdata_test.go` | Tests for log data conversion, attribute counts, severity | High |
| `config_test.go` | Tests for config loading from YAML | Medium |
| `factory_test.go` | Tests for factory creation, default config, K8s client mock | Medium |
| `testdata/config.yaml` | Test config fixtures | Medium — update when adding config fields |

---

## KEY DEPENDENCIES

| Dependency | Import Path | Purpose |
|------------|-------------|---------|
| k8sconfig | `internal/k8sconfig` | Shared K8s auth config (`APIConfig`, `AuthType`, `MakeClient`, `MakeDynamicClient`) |
| k8sinventory | `internal/k8sinventory` | Shared `GetEventTimestamp()`, `Config` for watch observer |
| watchobserver | `internal/k8sinventory/watch` | `Observer` — manages K8s API watches per namespace with retry, replaces old informer pattern |
| k8sleaderelector | `extension/k8sleaderelector` | `LeaderElection` interface for HA leader election integration |
| k8sleaderelectortest | `internal/k8sleaderelectortest` | `FakeLeaderElection`, `FakeHost` for leader election test mocking |
| dynamic | `k8s.io/client-go/dynamic` | Dynamic K8s API client for unstructured watches |
| dynamicfake | `k8s.io/client-go/dynamic/fake` | `NewSimpleDynamicClient()` for unit tests |
| client-go | `k8s.io/client-go` | K8s API client, fake client for tests |
| core/v1 | `k8s.io/api/core/v1` | `Event`, `ObjectReference`, `EventSource` structs |
| unstructured | `k8s.io/apimachinery/pkg/apis/meta/v1/unstructured` | Dynamic client returns `Unstructured` objects |
| runtime | `k8s.io/apimachinery/pkg/runtime` | `DefaultUnstructuredConverter` for Unstructured → typed conversion |
| schema | `k8s.io/apimachinery/pkg/runtime/schema` | `GroupVersionResource` for events GVR definition |
| apiWatch | `k8s.io/apimachinery/pkg/watch` | `EventType` constants (Added, Deleted, etc.) |
| plog | `go.opentelemetry.io/collector/pdata/plog` | OTel log data model |
| pcommon | `go.opentelemetry.io/collector/pdata/pcommon` | Timestamps, attribute values |
| semconv | `go.opentelemetry.io/otel/semconv/v1.38.0` | Semantic conventions (`K8SNodeNameKey`, `K8SNamespaceNameKey`) |
| receiverhelper | `go.opentelemetry.io/collector/receiver/receiverhelper` | `ObsReport` for instrumented receive operations |
| consumertest | `go.opentelemetry.io/collector/consumer/consumertest` | `LogsSink`, `NewNop()`, `NewErr()` for test consumers |
| fake | `k8s.io/client-go/kubernetes/fake` | `NewClientset()` for unit tests (typed client) |
| goleak | `go.uber.org/goleak` | Goroutine leak detection in `TestMain` |

---

## COMPONENT LIFECYCLE

### Startup Sequence
```
NewFactory()
  └── createLogsReceiver(ctx, settings, cfg, consumer)
        └── newReceiver(settings, config, consumer)
              ├── Creates ObsReport
              └── Returns &k8seventsReceiver{startTime: time.Now()}

Start(ctx, host)
  ├── context.WithCancel(ctx) → kr.ctx, kr.cancel
  ├── config.getDynamicClient() → dynamic.Interface
  │
  ├── IF K8sLeaderElector configured:
  │   ├── Resolve extension from host.GetExtensions()
  │   ├── Cast to k8sleaderelector.LeaderElection
  │   └── Register callbacks:
  │       ├── OnLeading:  create new ctx/cancel, call startWatchers()
  │       └── OnStopping: call Shutdown()
  │
  └── ELSE (no leader election):
      └── startWatchers()
            ├── Define eventsGVR = {Group: "", Version: "v1", Resource: "events"}
            ├── Build namespace list (empty string = all namespaces)
            ├── Create watchobserver.Observer with:
            │   ├── IncludeInitialState: false (don't replay historical events)
            │   ├── Exclude: {Deleted: true} (skip event deletions)
            │   └── Callback: Unstructured → corev1.Event → handleEvent()
            └── observer.Start(ctx, &wg) → returns stopper channel
                └── Mutex-protected: append stopper to stopperChanList
```

### Shutdown Sequence
```
Shutdown(ctx)
  ├── kr.cancel() (if cancel != nil) → cancels receiver context
  ├── Mutex-protected:
  │   ├── Close all stopper channels → stops watch goroutines
  │   └── Set stopperChanList = nil (allows re-start with leader election)
  └── kr.wg.Wait() → waits for all goroutines to finish
```

### Leader Election Lifecycle
When `k8s_leader_elector` is configured, the receiver supports stop/restart cycles:
```
Start() → registers callbacks, does NOT start watching
  OnLeading()  → creates new ctx/cancel, calls startWatchers()
  OnStopping() → calls Shutdown(), which closes channels + waits
  OnLeading()  → re-creates ctx/cancel, starts watching again (stopperChanList was nil'd)
```
This allows the receiver to stop processing when leadership is lost and resume when regained, without needing to recreate the receiver.

### Critical invariant: `startTime` guard
The receiver sets `startTime = time.Now()` at construction. `allowEvent()` rejects any event whose timestamp is before `startTime`. This prevents flooding the pipeline with historical events on startup.

---

## EVENT FILTERING PIPELINE (`allowEvent`)

```
allowEvent(ev) → (attributes []KeyValue, allow bool)
  │
  ├── 1. Timestamp guard: reject if getEventTimestamp(ev) < startTime
  │
  ├── 2. EventType filter (if configured):
  │   └── Reject if ev.Type not in config.EventTypes
  │
  ├── 3. ExcludeNamespaces gate (O(1) map lookup on excludedNSSet):
  │   ├── Rule 3 (Namespaces empty — watching all): primary filter,
  │   │   drops events whose InvolvedObject.Namespace is in excludedNSSet
  │   ├── Rule 4 (both Namespaces + ExcludeNamespaces non-empty):
  │   │   excluded namespaces already removed in startWatchers() set-difference;
  │   │   gate is a safe no-op here
  │   └── excludedNSSet is nil when ExcludeNamespaces is empty → nil map lookup is safe in Go
  │
  └── 4. IncludeInvolvedObject filter (if configured):
      ├── Look up ev.InvolvedObject.Kind in config map
      ├── If not found, try "Other" fallback key
      ├── If neither found → reject
      └── If found with IncludeReasons:
          ├── Look up ev.Reason in reason list
          ├── If not found → reject
          └── If found → return reason's custom Attributes
```

### Namespace Filtering Rules

| Rule | `namespaces` | `exclude_namespaces` | Behaviour |
|------|-------------|----------------------|-----------|
| 1 | empty | empty | Watch all namespaces; no exclusion |
| 2 | non-empty | empty | Watch listed namespaces only |
| 3 | empty | non-empty | Watch all namespaces; `allowEvent()` drops excluded per-event |
| 4 | non-empty | non-empty | Set-difference in `startWatchers()` removes excluded before creating watches; `allowEvent()` gate is a no-op |

**Rule 4 guard:** if set-difference produces an empty list (all namespaces excluded), `startWatchers()` logs a Warn and returns early — no watch connections are created.

### Filter config example
```yaml
k8s_events:
  namespaces: [default, kube-system]
  exclude_namespaces: [kube-system]   # Rule 4: set-difference → watches only [default]
  event_types: [Warning]
  include_involved_objects:
    Pod:
      include_reasons:
        - name: OOMKilled
          attributes:
            - key: alert.severity
              value: critical
        - name: CrashLoopBackOff
    Node:
      include_reasons:
        - name: NotReady
    Other:
      include_reasons:
        - name: FailedScheduling
```

---

## EVENT-TO-LOG DATA MODEL

### Resource Attributes
| Attribute | Source | Notes |
|-----------|--------|-------|
| `k8s.node.name` | `ev.ReportingInstance` or `ev.Source.Host` | Only set if non-empty |
| `k8s.object.kind` | `ev.InvolvedObject.Kind` | Always set |
| `k8s.object.name` | `ev.InvolvedObject.Name` | Always set |
| `k8s.object.uid` | `ev.InvolvedObject.UID` | Always set |
| `k8s.object.fieldpath` | `ev.InvolvedObject.FieldPath` | Always set |
| `k8s.object.api_version` | `ev.InvolvedObject.APIVersion` | Always set |
| `k8s.object.resource_version` | `ev.InvolvedObject.ResourceVersion` | Always set |
| `k8s.namespace.name` | `ev.InvolvedObject.Namespace` | Only if non-empty |
| `k8s.<kind>.name` | `ev.InvolvedObject.Name` | Dynamic key: e.g. `k8s.pod.name`, `k8s.node.name` |
| `resourceName` | `ev.InvolvedObject.Name` | Duplicate for compatibility |
| `type` | hardcoded `"event"` | TODO: should come from config |

### Log Record Fields
| Field | Source |
|-------|--------|
| Timestamp | `k8sinventory.GetEventTimestamp(ev)` — priority: EventTime > LastTimestamp > FirstTimestamp |
| Body | `ev.Message` |
| SeverityNumber | `"normal"` → Info, `"warning"` → Warn, `"error"` → Error, `"critical"` → Fatal (case-insensitive lookup) |
| SeverityText | `ev.Type` (raw) |
| Scope.Name | `metadata.ScopeName` |
| Scope.Version | `settings.BuildInfo.Version` |

### Log Attributes
| Attribute | Source | Notes |
|-----------|--------|-------|
| `k8s.event.type` | `ev.Type` | "Normal" or "Warning" |
| `k8s.event.sourceComponent` | `ev.Source.Component` | |
| `k8s.event.reason` | `ev.Reason` | |
| `k8s.event.action` | `ev.Action` | |
| `k8s.event.start_time` | `ev.ObjectMeta.CreationTimestamp` | String format |
| `k8s.event.name` | `ev.ObjectMeta.Name` | |
| `k8s.event.uid` | `ev.ObjectMeta.UID` | |
| `level` | `ev.Type` | Duplicate for compatibility |
| `k8s.namespace.name` | `ev.InvolvedObject.Namespace` | Always set (even if empty) |
| `k8s.event.count` | `ev.Count` | Only set if Count != 0 |
| *(custom)* | `attributes` from `allowEvent()` | From `ReasonProperties.Attributes` config |

---

## REQUIRED PATTERNS

### R1 — ObsReport Bracketing for ConsumeLogs
Every call to `logsConsumer.ConsumeLogs()` MUST be wrapped in `obsrecv.StartLogsOp` / `obsrecv.EndLogsOp`:
```go
ctx := kr.obsrecv.StartLogsOp(kr.ctx)
consumerErr := kr.logsConsumer.ConsumeLogs(ctx, ld)
kr.obsrecv.EndLogsOp(ctx, metadata.Type.String(), 1, consumerErr)
```

### R2 — Test Mocking via makeClient and makeDynamicClient
Tests MUST override BOTH `config.makeClient` and `config.makeDynamicClient`:
```go
rCfg.makeClient = func(k8sconfig.APIConfig) (k8s.Interface, error) {
    return fake.NewClientset(), nil
}
scheme := runtime.NewScheme()
_ = corev1.AddToScheme(scheme)
rCfg.makeDynamicClient = func(k8sconfig.APIConfig) (dynamic.Interface, error) {
    return dynamicfake.NewSimpleDynamicClient(scheme), nil
}
```
The dynamic client requires a scheme with `corev1` registered for unstructured conversion to work. Never use real K8s clients in unit tests.

### R3 — Goroutine Leak Detection
The file `generated_package_test.go` runs `goleak.VerifyTestMain(m)`. All tests MUST properly shut down goroutines. If adding informer goroutines, ensure stopper channels are closed in `Shutdown()`.

### R4 — Event Timestamp Priority
Always respect: `EventTime > LastTimestamp > FirstTimestamp`. The `getEventTimestamp()` function encodes this. Do not change the priority without understanding K8s event deprecation strategy (EventTime is the modern field).

### R5 — Stopper Channel Lifecycle with Mutex and WaitGroup
Stopper channels MUST be protected by `kr.mu` and goroutines tracked by `kr.wg`:
```go
// In startWatchers:
stopperChan := observer.Start(kr.ctx, &kr.wg)
kr.mu.Lock()
kr.stopperChanList = append(kr.stopperChanList, stopperChan)
kr.mu.Unlock()

// In Shutdown:
kr.mu.Lock()
for _, stopperChan := range kr.stopperChanList {
    close(stopperChan)
}
kr.stopperChanList = nil  // allows re-start with leader election
kr.mu.Unlock()
kr.wg.Wait()  // wait for all goroutines to finish
```

### R6 — Config Validation
The `Validate()` method MUST check every new config field if and only if required, else skip this. Event types are validated against `EventTypeNormal` / `EventTypeWarning` constants. Always call `cfg.APIConfig.Validate()` at the end.

### R7 — Attribute Capacity Hints
When adding new attributes to `k8sEventToLogData`, update the capacity constants:
```go
const (
    totalLogAttributes      = 7  // ← update thisgi
    totalResourceAttributes = 6  // ← update this
)
```
And call `attrs.EnsureCapacity(totalLogAttributes)` / `resourceAttrs.EnsureCapacity(totalResourceAttributes)`.

### R12 — ExcludeNamespaces: Use toSet() and Build at Construction
When working with `ExcludeNamespaces` or any string-set lookup, always use the `toSet()` helper — never inline the `make` + range loop:
```go
// ✅ REQUIRED
excludedNSSet = toSet(config.ExcludeNamespaces)

// ❌ FORBIDDEN: duplicate inline set-building
excludedNSSet = make(map[string]struct{}, len(config.ExcludeNamespaces))
for _, ns := range config.ExcludeNamespaces {
    excludedNSSet[ns] = struct{}{}
}
```
The `excludedNSSet` field MUST be built in `newReceiver()` (not `startWatchers()`), so it is available for `allowEvent()` in tests that create a receiver directly without calling `Start()`.

### R13 — startWatchers() Rule 4 Empty-Set Guard
After computing the namespace set-difference, always guard against the all-excluded case:
```go
if len(namespaces) > 0 && len(kr.config.ExcludeNamespaces) > 0 {
    namespaces = computeNamespaceFilter(namespaces, kr.config.ExcludeNamespaces)
    if len(namespaces) == 0 {
        kr.settings.Logger.Warn("all namespaces excluded by exclude_namespaces filter — receiver will not collect any events")
        return  // ← REQUIRED: without this, len==0 falls through to watch-all
    }
}
```
Omitting this guard causes the opposite of the intended behaviour: an empty list is treated as "watch all namespaces".

### R8 — Changelog Entry
Every PR that changes behavior MUST add a `.yaml` file to `.chloggen/`. Use `make chlog-new` to generate. Required for: config changes, attribute changes, filter logic changes, bug fixes.

### R9 — Leader Election Integration
When adding or modifying leader election support:
- The `K8sLeaderElector` config field is an optional `*component.ID`
- Extension is resolved via `host.GetExtensions()` and cast to `k8sleaderelector.LeaderElection`
- `SetCallBackFuncs(onLeading, onStopping)` registers lifecycle callbacks
- `onLeading` MUST create a new `context.WithCancel` and call `startWatchers()`
- `onStopping` MUST call `Shutdown()` to clean up
- Callbacks remain registered for the receiver lifetime — re-election restarts watching
```go
elector.SetCallBackFuncs(
    func(ctx context.Context) {
        cctx, cancel := context.WithCancel(ctx)
        kr.cancel = cancel
        kr.ctx = cctx
        kr.startWatchers()
    },
    func() {
        kr.Shutdown(context.Background())
    })
```

### R10 — Dynamic Client and Watch Observer
The receiver uses `dynamic.Interface` (not typed `k8s.Interface`) for watching events:
- Events are received as `*unstructured.Unstructured` objects
- MUST convert using `runtime.DefaultUnstructuredConverter.FromUnstructured()` to `*corev1.Event`
- Watch observer is created via `watchobserver.New()` with `IncludeInitialState: false` and `Exclude: {Deleted: true}`
- The observer handles namespace scoping, retry, and resource version management internally
```go
observer, err := watchobserver.New(
    kr.client,
    watchobserver.Config{
        Config: k8sinventory.Config{
            Gvr:        eventsGVR,
            Namespaces: namespaces,
        },
        IncludeInitialState: false,
        Exclude:             map[apiWatch.EventType]bool{apiWatch.Deleted: true},
    },
    kr.settings.Logger,
    func(event *apiWatch.Event) {
        // Convert unstructured to corev1.Event, then handleEvent()
    },
)
```

### R11 — Scope Instrumentation
Log data MUST set scope name and version:
```go
sl.Scope().SetName(metadata.ScopeName)
sl.Scope().SetVersion(version)  // from settings.BuildInfo.Version
```

---

## FORBIDDEN PATTERNS

### F1 — Blocking in Watch Observer Callback
```go
// ❌ FORBIDDEN: Blocking call inside the watch observer callback
func(event *apiWatch.Event) {
    time.Sleep(5 * time.Second)  // blocks the watch event loop
    kr.handleEvent(ev)
}
```
The watch observer callback runs in the observer's goroutine. Long-running work should be dispatched asynchronously or kept fast.

### F2 — Missing ObsReport Instrumentation
```go
// ❌ FORBIDDEN: ConsumeLogs without ObsReport
func (kr *k8seventsReceiver) handleEvent(ev *corev1.Event) {
    ld := k8sEventToLogData(kr.settings.Logger, ev, attributes)
    kr.logsConsumer.ConsumeLogs(kr.ctx, ld)  // No observability!
}
```
```go
// ✅ REQUIRED: Always bracket with ObsReport
ctx := kr.obsrecv.StartLogsOp(kr.ctx)
consumerErr := kr.logsConsumer.ConsumeLogs(ctx, ld)
kr.obsrecv.EndLogsOp(ctx, metadata.Type.String(), 1, consumerErr)
```

### F3 — Bypassing the startTime Guard
```go
// ❌ FORBIDDEN: Processing events without timestamp check
func (kr *k8seventsReceiver) handleEvent(ev *corev1.Event) {
    ld := k8sEventToLogData(kr.settings.Logger, ev, nil)
    kr.logsConsumer.ConsumeLogs(kr.ctx, ld)
}
```
Always go through `allowEvent()` which enforces the startup timestamp guard.

### F4 — Editing Generated Files
```go
// ❌ FORBIDDEN: Manual edits to these files
// - internal/metadata/generated_status.go
// - generated_component_test.go
// - generated_package_test.go
```
These are generated by `mdatagen` from `metadata.yaml`. Edit `metadata.yaml` and run `go generate ./...` instead.

### F5 — Real K8s Client in Tests
```go
// ❌ FORBIDDEN: Using real cluster in unit tests
client, _ := kubernetes.NewForConfig(config)
```
```go
// ✅ REQUIRED: Use fake clients
client := fake.NewClientset()
scheme := runtime.NewScheme()
_ = corev1.AddToScheme(scheme)
dynClient := dynamicfake.NewSimpleDynamicClient(scheme)
```

### F6 — Hardcoded Severity Values
```go
// ❌ FORBIDDEN: Hardcoding severity numbers
lr.SetSeverityNumber(9)  // Magic number
```
```go
// ✅ REQUIRED: Use plog constants via severityMap
if severityNumber, ok := severityMap[strings.ToLower(ev.Type)]; ok {
    lr.SetSeverityNumber(severityNumber)
    lr.SetSeverityText(ev.Type)
}
```

### F7 — Using client-go Informers (Deprecated Pattern)
```go
// ❌ FORBIDDEN: Old informer-based watching (replaced by watch observer)
cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "events", ns, fields.Everything())
cache.NewInformerWithOptions(...)
```
```go
// ✅ REQUIRED: Use k8sinventory/watch.Observer with dynamic client
observer, err := watchobserver.New(kr.client, config, logger, callback)
stopperChan := observer.Start(kr.ctx, &kr.wg)
```

### F9 — Duplicating Set-Building Logic
```go
// ❌ FORBIDDEN: Inline set-building anywhere in the package
excluded := make(map[string]struct{}, len(items))
for _, item := range items {
    excluded[item] = struct{}{}
}
```
```go
// ✅ REQUIRED: Use the shared toSet() helper
excluded := toSet(items)
```
The `toSet()` function exists precisely to avoid this duplication. It is used in both `newReceiver()` and `computeNamespaceFilter()`.

### F10 — Building excludedNSSet in startWatchers()
```go
// ❌ FORBIDDEN: building the map inside startWatchers()
func (kr *k8seventsReceiver) startWatchers() {
    if len(kr.config.ExcludeNamespaces) > 0 {
        kr.excludedNSSet = toSet(kr.config.ExcludeNamespaces)  // too late for tests
    }
    // ...
}
```
```go
// ✅ REQUIRED: build in newReceiver() so allowEvent() tests work without Start()
if len(config.ExcludeNamespaces) > 0 {
    excludedNSSet = toSet(config.ExcludeNamespaces)
}
```

### F8 — Modifying stopperChanList Without Mutex
```go
// ❌ FORBIDDEN: Unsynchronized access to stopperChanList
kr.stopperChanList = append(kr.stopperChanList, stopperChan)
```
```go
// ✅ REQUIRED: Mutex-protected access
kr.mu.Lock()
kr.stopperChanList = append(kr.stopperChanList, stopperChan)
kr.mu.Unlock()
```

---

## COMMON TASK RECIPES

### Adding a New Namespace Filter (reference: exclude_namespaces)

The `exclude_namespaces` field added in ITOM-116409 is the canonical example for namespace-level filtering:

1. Config field: `ExcludeNamespaces []string \`mapstructure:"exclude_namespaces,omitempty"\``
2. Receiver struct field: `excludedNSSet map[string]struct{}` (built with `toSet()` in `newReceiver()`)
3. startWatchers() Rule 4 block with empty-set guard (see R13)
4. allowEvent() gate using `kr.excludedNSSet` map lookup (see R12 + filter pipeline above)
5. Tests: `TestComputeNamespaceFilter` (6 cases) + `TestAllowEventExcludeNamespaces` (8 cases incl. cluster-scoped empty-namespace events)

---

### Adding a New Config Field

1. **Add field to `Config` struct** in `config.go`:
   ```go
   type Config struct {
       k8sconfig.APIConfig `mapstructure:",squash"`
       Namespaces          []string    `mapstructure:"namespaces"`
       ExcludeNamespaces   []string    `mapstructure:"exclude_namespaces,omitempty"`
       EventTypes          []EventType `mapstructure:"event_types,omitempty"`
       NewField            string      `mapstructure:"new_field,omitempty"`  // ← add here
       // ...
   }
   ```

2. **Add validation** in `config.go` `Validate()` method:
   ```go
   func (cfg *Config) Validate() error {
       if cfg.NewField != "" {
           // validate...
       }
       // existing validations...
       return cfg.APIConfig.Validate()
   }
   ```

3. **Add test config** in `testdata/config.yaml`:
   ```yaml
   k8s_events/all_settings:
     namespaces: [default, my_namespace]
     new_field: "value"
   ```

4. **Add config test case** in `config_test.go`:
   ```go
   {
       id: component.NewIDWithName(metadata.Type, "all_settings"),
       expected: &Config{
           Namespaces: []string{"default", "my_namespace"},
           NewField:   "value",
           APIConfig:  k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeServiceAccount},
       },
   },
   ```

5. **Update default config** in `factory.go` `createDefaultConfig()` if field has a non-zero default.

6. **Update README.md** with the new config option.

7. **Add changelog entry**: `make chlog-new`

---

### Adding a New Log/Resource Attribute

1. **Add attribute** in `k8sEventToLogData()` in `k8s_event_to_logdata.go`:
   ```go
   // For resource attributes:
   resourceAttrs.PutStr("k8s.new.attribute", ev.SomeField)

   // For log attributes:
   attrs.PutStr("k8s.event.new_attr", ev.SomeField)
   ```

2. **Update capacity constants**:
   ```go
   const (
       totalLogAttributes      = 8  // was 7, now 8
       totalResourceAttributes = 6  // update if resource attribute
   )
   ```

3. **Update test assertions** in `k8s_event_to_logdata_test.go`:
   ```go
   assert.Equal(t, 8, attrs.Len())  // updated count
   ```

4. **Add specific attribute test**:
   ```go
   attr, ok := attrs.Get("k8s.event.new_attr")
   assert.True(t, ok)
   assert.Equal(t, "expected_value", attr.AsString())
   ```

---

### Adding a New Event Filter

1. **Define filter types** in `config.go` (if needed):
   ```go
   type NewFilterProp struct {
       Field string `mapstructure:"field"`
   }
   ```

2. **Add filter field to Config**:
   ```go
   type Config struct {
       // ...existing fields...
       NewFilter []NewFilterProp `mapstructure:"new_filter,omitempty"`
   }
   ```

3. **Add filter logic** to `allowEvent()` in `receiver.go` — add AFTER the timestamp guard:
   ```go
   func (kr *k8seventsReceiver) allowEvent(ev *corev1.Event) (attributes []KeyValue, allow bool) {
       eventTimestamp := getEventTimestamp(ev)
       if eventTimestamp.Before(kr.startTime) {
           return attributes, false
       }
       // ...existing filters...

       // New filter — add here
       if len(kr.config.NewFilter) != 0 {
           // filter logic...
       }

       return attributes, true
   }
   ```

4. **Add tests** in `receiver_test.go`:
   ```go
   func TestAllowEventWithNewFilter(t *testing.T) {
       rCfg := createDefaultConfig().(*Config)
       rCfg.NewFilter = []NewFilterProp{{Field: "value"}}
       r, err := newReceiver(receivertest.NewNopSettings(metadata.Type), rCfg, consumertest.NewNop())
       require.NoError(t, err)
       recv := r.(*k8seventsReceiver)

       k8sEvent := getEvent("Normal")
       _, allow := recv.allowEvent(k8sEvent)
       // assert...
   }
   ```

---

### Adding Support for a New K8s Event Field

When K8s adds new fields to `corev1.Event`, to surface them:

1. Check if the `k8s.io/api` dependency needs upgrading in `go.mod`
2. Map the field to either resource attribute or log attribute (see data model above)
3. Follow the "Adding a New Log/Resource Attribute" recipe
4. If the field affects filtering, also follow the "Adding a New Event Filter" recipe

---

### Fixing a Bug in Event Processing

1. **Write a failing test first** using `getEvent(eventType)` helper:
   ```go
   func TestBugDescription(t *testing.T) {
       rCfg := createDefaultConfig().(*Config)
       sink := new(consumertest.LogsSink)
       r, err := newReceiver(receivertest.NewNopSettings(metadata.Type), rCfg, sink)
       require.NoError(t, err)
       recv := r.(*k8seventsReceiver)
       recv.ctx = t.Context()

       k8sEvent := getEvent("Normal")
       // Set up event to trigger the bug
       k8sEvent.SomeField = "trigger_value"

       recv.handleEvent(k8sEvent)
       // Assert expected behavior
       assert.Equal(t, 1, sink.LogRecordCount())
   }
   ```

2. **Fix the code**

3. **Verify test passes**

4. **Check for goroutine leaks**: `go test -run TestBugDescription -count=1 ./...`

---

## TESTING CONVENTIONS

### Test Structure
```
receiver_test.go          — Lifecycle, handleEvent, allowEvent, getEventTimestamp,
                            leader election, dynamic client error paths,
                            multiple namespace watching, consumer error handling
k8s_event_to_logdata_test.go — Log data conversion, attribute mapping, severity
                                (Normal/Warning/Error/Critical/Unknown),
                                scope name/version, api/resource version
config_test.go             — Config loading from YAML testdata
factory_test.go            — Factory creation, default config, receiver creation
                              with both typed and dynamic fake clients
generated_component_test.go — (AUTO) Component lifecycle
generated_package_test.go   — (AUTO) goleak.VerifyTestMain
```

### Test Helpers
- `getEvent(eventType string)` — Returns a well-formed `*corev1.Event` with all required fields populated (Pod, with timestamps, source, etc.). Takes `eventType` parameter ("Normal", "Warning", etc.)
- `fake.NewClientset()` — Mock typed K8s client
- `dynamicfake.NewSimpleDynamicClient(scheme)` — Mock dynamic K8s client (requires `corev1.AddToScheme(scheme)`)
- `consumertest.LogsSink` — Captures consumed logs for assertion
- `consumertest.NewNop()` — No-op consumer when logs aren't inspected
- `consumertest.NewErr(err)` — Error-returning consumer for error path testing
- `receivertest.NewNopSettings(metadata.Type)` — Default receiver settings with component type
- `componenttest.NewNopHost()` — No-op host for Start()
- `confmaptest.LoadConf(path)` — Loads YAML config for testing
- `k8sleaderelectortest.FakeLeaderElection` — Fake leader elector with `InvokeOnLeading()`/`InvokeOnStopping()` for lifecycle testing
- `k8sleaderelectortest.FakeHost` — Host that returns the fake leader elector extension
- `t.Context()` — Use test context (preferred over `context.Background()` in tests)

### Running Tests
```bash
cd receiver/k8seventsreceiver
go test ./...                         # All tests
go test -run TestHandleEvent ./...    # Single test
go test -race ./...                   # Race detector
go test -count=1 ./...                # No cache
```

### When to Update `getEvent()` Helper
If your change adds a new field that MUST be present on typical events, add it to `getEvent()`. Keep it minimal — only fields that most tests need.

---

## BUILD & DEVELOPMENT

### Module Structure
- Module: `github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver`
- Key dependencies: `internal/k8sconfig`, `internal/k8sinventory`, `internal/k8sleaderelectortest`, `extension/k8sleaderelector`
- Build: `make` (uses `../../Makefile.Common`)

### Code Generation
After editing `metadata.yaml`:
```bash
go generate ./...
```
This regenerates `internal/metadata/generated_status.go`.

### Linting
```bash
make lint          # From receiver/k8seventsreceiver/
```

---

## KNOWN QUIRKS & GOTCHAS

1. **`allowEvent` returns two values** — `(attributes []KeyValue, allow bool)`. The `attributes` come from `ReasonProperties.Attributes` in config. Both values must be checked; attributes are passed to `k8sEventToLogData`.

2. **Watch observer `IncludeInitialState: false`** — The observer does NOT replay existing events on startup. Only new watch events are processed. This is intentional to avoid duplicate processing.

3. **Deleted events are excluded** — The watch observer is configured with `Exclude: {Deleted: true}`. K8s event deletions are not processed (they are garbage collected, not meaningful).

4. **`startTime` is set at `newReceiver()`, not `Start()`** — This means the flood guard uses construction time, not start time. This is a subtle but intentional choice.

5. **`IncludeInvolvedObject` "Other" fallback** — If an event's `InvolvedObject.Kind` is not in the config map, the code looks for an "Other" key as a catch-all before rejecting.

6. **Dynamic resource attribute key** — `k8s.<kind>.name` is constructed dynamically from `strings.ToLower(ev.InvolvedObject.Kind)`. This means every object kind gets a different attribute key.

7. **Severity mapping is case-insensitive and extended** — `strings.ToLower(ev.Type)` is used to look up in `severityMap`. The map includes 4 levels: normal→Info, warning→Warn, error→Error, critical→Fatal.

8. **`type: "event"` is hardcoded** — There's a TODO in the code to make this configurable.

9. **`excludedNSSet` nil map lookup is safe** — In Go, a nil map lookup `m[key]` returns the zero value and does not panic. When `ExcludeNamespaces` is empty, `excludedNSSet` is nil and the `allowEvent()` gate `if _, excluded := kr.excludedNSSet[ns]; excluded` correctly evaluates to `false` without any nil guard.

10. **`computeNamespaceFilter` returns nil (not empty slice) when all excluded** — The `filtered` var is declared as `var filtered []string` (nil). Appending nothing leaves it nil. `len(nil) == 0`, so the caller's `if len(namespaces) == 0` guard correctly triggers the early return.

9. **`k8sEventToLogData` takes 4 args** — `(logger, event, version, attributes)`. The `version` parameter sets `Scope.Version` from `settings.BuildInfo.Version`. Update tests when modifying the function signature.

10. **Dynamic client returns Unstructured objects** — The watch observer uses `dynamic.Interface`, so events arrive as `*unstructured.Unstructured`. The receiver callback MUST convert them to `*corev1.Event` using `runtime.DefaultUnstructuredConverter.FromUnstructured()`.

11. **Leader election allows stop/restart cycles** — When configured with `k8s_leader_elector`, `Shutdown()` sets `stopperChanList = nil`, and the `onLeading` callback creates a fresh context. This allows the receiver to restart without recreation.

12. **`getEventTimestamp` is now in shared package** — The timestamp priority logic lives in `k8sinventory.GetEventTimestamp()`. A local `getEventTimestamp()` wrapper still exists in `receiver.go` for `allowEvent()` usage, but `k8sEventToLogData` calls the shared version directly.

13. **Semconv uses v1.38.0** — The import is `go.opentelemetry.io/otel/semconv/v1.38.0`, using `conventions.K8SNodeNameKey` and `conventions.K8SNamespaceNameKey` (not the old collector semconv package).

---

## ENHANCEMENT IDEAS & KNOWN TODOS

These are areas the maintainers have flagged or that are natural next steps:

- **Make `type` resource attribute configurable** (currently hardcoded to `"event"`)
- **Promote from alpha to beta** — requires passing lifecycle tests (currently skipped: `skip_lifecycle: true` in `metadata.yaml`)
- **Events.v1 API support** — K8s is moving from `core/v1.Event` to `events.k8s.io/v1.Event`
- **ExcludeInvolvedObject** filter (inverse of IncludeInvolvedObject)
- **Regex-based filtering** for object names, reasons, messages
- **Configurable attribute mapping** — let users choose which event fields become which attributes
- **Batch event delivery** — currently sends one log per event; batching could improve throughput
- **Update capacity constants** — `totalLogAttributes` (7) and `totalResourceAttributes` (6) may be stale after attribute additions
