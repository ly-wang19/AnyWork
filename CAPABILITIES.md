# AnyWork capabilities

[简体中文](CAPABILITIES.zh-CN.md) · **English** · [日本語](CAPABILITIES.ja.md)

AnyWork gives an AI agent repeatable methods for turning ambiguous work into checked, usable deliverables. The current alpha contains **24 atomic Skills** and **7 installable work packs**. Skills can be used alone or composed into an end-to-end workflow.

## Capability map

| Work area | What you can get | Included Skills |
|---|---|---|
| Define and plan | A testable task contract, executable plan, or ranked backlog | [Clarify outcome](skills/clarify-outcome/SKILL.md), [plan work](skills/plan-work/SKILL.md), [prioritize work](skills/prioritize-work/SKILL.md) |
| Research and evidence | A traceable source log, structured extraction, claim verdicts, or evidence synthesis | [Gather sources](skills/gather-sources/SKILL.md), [extract structured information](skills/extract-structured-info/SKILL.md), [verify claims](skills/verify-claims/SKILL.md), [synthesize sources](skills/synthesize-sources/SKILL.md) |
| Analyze and decide | Reproducible findings, a root-cause report, transparent scenarios, an option comparison, or a decision recommendation | [Analyze data](skills/analyze-data/SKILL.md), [diagnose root cause](skills/diagnose-root-cause/SKILL.md), [model scenarios](skills/model-scenarios/SKILL.md), [compare options](skills/compare-options/SKILL.md), [make a decision](skills/make-decision/SKILL.md) |
| Create and communicate | A decision brief, long-form document, auditable workbook, presentation, ready-to-send message, or localized asset | [Write a brief](skills/write-brief/SKILL.md), [draft a document](skills/draft-document/SKILL.md), [build a spreadsheet](skills/build-spreadsheet/SKILL.md), [design a presentation](skills/design-presentation/SKILL.md), [write a message](skills/write-message/SKILL.md), [localize content](skills/localize-content/SKILL.md) |
| Meetings and operations | A meeting design, decision/action log, executable SOP, or safe automation plan | [Prepare a meeting](skills/prepare-meeting/SKILL.md), [capture decisions and actions](skills/capture-decisions-actions/SKILL.md), [create an SOP](skills/create-sop/SKILL.md), [automate a routine](skills/automate-routine/SKILL.md) |
| Quality and improvement | A deliverable review or evidence-based retrospective with owned experiments | [Review a deliverable](skills/review-deliverable/SKILL.md), [run a retrospective](skills/run-retrospective/SKILL.md) |

## Work packs

Packs are curated starting points. Their Skills intentionally overlap because the same atomic capability can support several roles.

| Pack ID | Designed for | Skills | Install example |
|---|---|---:|---|
| `essential` | Everyday research, writing, coordination, and review | 9 | `anywork install essential` |
| `manager-leadership` | Decisions, priorities, scenarios, meetings, and retrospectives | 8 | `anywork install manager-leadership` |
| `product-operations` | Product discovery, prioritization, delivery systems, and improvement | 9 | `anywork install product-operations` |
| `research-consulting` | Evidence-heavy research, analysis, recommendations, reports, and decks | 7 | `anywork install research-consulting` |
| `go-to-market` | Market evidence, customer decisions, content, communication, and execution | 11 | `anywork install go-to-market` |
| `people-recruiting` | Structured, fair, and auditable people and recruiting work | 10 | `anywork install people-recruiting` |
| `complete` | Every official atomic Skill in this release | 24 | `anywork install complete` |

Add `--agent codex`, `--agent claude-code`, `--agent gemini-cli`, `--agent github-copilot`, or `--agent opencode` to target one documented host. Use `--scope project` or `--scope user` to choose the installation boundary, and use `--dry-run` before writing.

## Example end-to-end workflows

### Research to recommendation

`clarify-outcome` → `gather-sources` → `verify-claims` → `synthesize-sources` → `compare-options` → `make-decision` → `write-brief`

Result: a decision-ready recommendation whose sources, uncertainty, alternatives, and rationale remain traceable.

### Data to executive presentation

`extract-structured-info` → `analyze-data` → `model-scenarios` → `build-spreadsheet` → `design-presentation` → `review-deliverable`

Result: a recalculable model and an audience-ready narrative, with the evidence and assumptions kept separate.

### Meeting to accountable execution

`prepare-meeting` → `capture-decisions-actions` → `plan-work` → `prioritize-work` → `review-deliverable`

Result: a focused meeting followed by explicit decisions, owners, deadlines, dependencies, and acceptance checks.

### Recurring operation to safe automation

`diagnose-root-cause` → `create-sop` → `automate-routine` → `run-retrospective`

Result: a controlled operating process with exceptions, human gates, recovery, and measured improvement.

### Global content delivery

`write-brief` → `draft-document` → `localize-content` → `write-message` → `review-deliverable`

Result: channel-ready English, Simplified Chinese, and Japanese content that preserves the approved meaning and terminology.

## Capability boundary

AnyWork installs work methods, not credentials or silent integrations. A Skill can use only the tools and access already available to its host agent. External writes, messages, purchases, deletions, and other consequential actions still require explicit authorization and a verifiable target. All 24 Skills remain `experimental` until the [quality gates](docs/QUALITY.md) and [roadmap](ROADMAP.md) evidence are complete.
