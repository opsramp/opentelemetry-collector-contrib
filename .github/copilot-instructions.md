# OpenTelemetry Collector Contrib — Copilot Instructions

## Skill Routing
Before responding, consult `.github/skills/skill-rules.json` to determine which skill best matches the user's request based on `promptTriggers` (keywords and intentPatterns), `fileTriggers` (path and content patterns).
- If a matching skill has `enforcement: "block"`, read and apply that skill's `SKILL.md` **before** proceeding.
- If a matching skill has `enforcement: "suggest"`, mention the skill is available but proceed normally.
- If multiple skills match, prefer the one with higher `priority`.
