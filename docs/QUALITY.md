# AnyWork quality standard

[简体中文](QUALITY.zh-CN.md) · **English** · [日本語](QUALITY.ja.md)

An official Skill is a releaseable work capability, not an interesting prompt.

## Required gates

1. Keep one atomic task and one canonical `SKILL.md`.
2. Provide English, Simplified Chinese, and Japanese trigger coverage.
3. Follow the user's output language and preserve source-language evidence.
4. Declare source, license, version, risk, capabilities, and side effects in the registry.
5. Avoid install-time execution, hidden network calls, telemetry, or secret collection.
6. Protect external writes with explicit authorization and verifiable targets.
7. Pass schema, installer, trigger, task-quality, safety, and cross-agent checks.
8. Preserve original attribution and licenses for adapted work.

## Stable-release targets

- Schema and render checks: 100%.
- Critical safety cases: 100%.
- Same-language response adherence: at least 98%.
- Trigger recall: at least 90%; false-positive rate: at most 5%.
- Average task score in every language: at least 85/100.
- Cross-language score gap: at most 5 percentage points.
- Deterministic script fixtures: 100%.
- Install lifecycle: macOS, Linux, and Windows × Codex and Claude Code × three locales.

Experimental Skills may enter the catalog before meeting the stable score thresholds, but they must pass structural, provenance, and critical safety gates and must be labeled `experimental`.

## Scoring rubric

- Core task correctness: 40%
- Evidence, traceability, and reproducibility: 20%
- Direct usability: 15%
- Native-language quality and register: 15%
- Safety, permission discipline, and non-fabrication: 10%
