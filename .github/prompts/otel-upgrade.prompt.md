---
description: "Upgrade opentelemetry-collector-contrib from one release branch to another with OpsRamp custom change preservation"
mode: "agent"
agent: "otel-upgrade"
---

Upgrade opentelemetry-collector-contrib from `${input:currentBranch}` to `${input:desiredBranch}`.

Follow all phases: branch creation & merge, build verification, OpsRamp package updates, and impact analysis.
