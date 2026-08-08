# AnyWork capabilities

[简体中文](CAPABILITIES.zh-CN.md) · **English** · [日本語](CAPABILITIES.ja.md)

AnyWork gives an AI agent repeatable methods for turning ambiguous work into checked, usable deliverables. The current alpha contains **24 atomic Skills, 1 end-to-end orchestrator, and 7 installable work packs**. Atomic Skills can be used alone; the orchestrator selects and executes the smallest useful chain for a complex outcome.

## Capability map

| Work area | What you can get | Included Skills |
|---|---|---|
| Define and plan | A testable task contract, executable plan, or ranked backlog | [Clarify outcome](skills/clarify-outcome/SKILL.md), [plan work](skills/plan-work/SKILL.md), [prioritize work](skills/prioritize-work/SKILL.md) |
| Research and evidence | A traceable source log, structured extraction, claim verdicts, or evidence synthesis | [Gather sources](skills/gather-sources/SKILL.md), [extract structured information](skills/extract-structured-info/SKILL.md), [verify claims](skills/verify-claims/SKILL.md), [synthesize sources](skills/synthesize-sources/SKILL.md) |
| Analyze and decide | Reproducible findings, a root-cause report, transparent scenarios, an option comparison, or a decision recommendation | [Analyze data](skills/analyze-data/SKILL.md), [diagnose root cause](skills/diagnose-root-cause/SKILL.md), [model scenarios](skills/model-scenarios/SKILL.md), [compare options](skills/compare-options/SKILL.md), [make a decision](skills/make-decision/SKILL.md) |
| Create and communicate | A decision brief, long-form document, auditable workbook, presentation, ready-to-send message, or localized asset | [Write a brief](skills/write-brief/SKILL.md), [draft a document](skills/draft-document/SKILL.md), [build a spreadsheet](skills/build-spreadsheet/SKILL.md), [design a presentation](skills/design-presentation/SKILL.md), [write a message](skills/write-message/SKILL.md), [localize content](skills/localize-content/SKILL.md) |
| Meetings and operations | A meeting design, decision/action log, executable SOP, or safe automation plan | [Prepare a meeting](skills/prepare-meeting/SKILL.md), [capture decisions and actions](skills/capture-decisions-actions/SKILL.md), [create an SOP](skills/create-sop/SKILL.md), [automate a routine](skills/automate-routine/SKILL.md) |
| Quality and improvement | A deliverable review or evidence-based retrospective with owned experiments | [Review a deliverable](skills/review-deliverable/SKILL.md), [run a retrospective](skills/run-retrospective/SKILL.md) |
| End-to-end orchestration | A checked multi-stage outcome with explicit handoffs, evidence, uncertainty, and permission gates | [Orchestrate work](skills/orchestrate-work/SKILL.md) |

## Work packs

Packs are curated starting points. Their Skills intentionally overlap because the same atomic capability can support several roles.

| Pack ID | Designed for | Skills | Install example |
|---|---|---:|---|
| `essential` | Everyday research, writing, coordination, review, and orchestration | 10 | `anywork install essential` |
| `manager-leadership` | Essential plus decisions, priorities, scenarios, meetings, and retrospectives | 18 | `anywork install manager-leadership` |
| `product-operations` | Essential plus product discovery, prioritization, delivery systems, and improvement | 19 | `anywork install product-operations` |
| `research-consulting` | Essential plus evidence-heavy analysis, recommendations, reports, and decks | 17 | `anywork install research-consulting` |
| `go-to-market` | Essential plus market evidence, customer decisions, content, communication, and execution | 21 | `anywork install go-to-market` |
| `people-recruiting` | Essential plus structured, fair, and auditable people and recruiting work | 20 | `anywork install people-recruiting` |
| `complete` | Every official atomic Skill and the orchestrator in this release | 25 | `anywork install complete` |

Add `--agent codex`, `--agent claude-code`, `--agent gemini-cli`, `--agent github-copilot`, or `--agent opencode` to target one documented host. Use `--scope project` or `--scope user` to choose the installation boundary, and use `--dry-run` before writing.

`anywork setup` installs `complete` for all documented hosts in one step. `anywork demo research`, `anywork demo meeting`, and `anywork demo automation` print ready-to-paste trilingual jobs that exercise the orchestrator.

## Example end-to-end workflows

### Research to recommendation

`orchestrate-work` → `clarify-outcome` when needed → `gather-sources` → `verify-claims` → `synthesize-sources` → `compare-options` → `make-decision` → `write-brief` → `review-deliverable`

Result: a decision-ready recommendation whose sources, uncertainty, alternatives, and rationale remain traceable.

### Data to executive presentation

`orchestrate-work` → `extract-structured-info` → `analyze-data` → `model-scenarios` → `build-spreadsheet` → `design-presentation` → `review-deliverable`

Result: a recalculable model and an audience-ready narrative, with the evidence and assumptions kept separate.

### Meeting to accountable execution

`orchestrate-work` → `prepare-meeting` → meeting occurs → `capture-decisions-actions` → `plan-work` → `prioritize-work` → `review-deliverable`

Result: a focused meeting followed by explicit decisions, owners, deadlines, dependencies, and acceptance checks.

### Recurring operation to safe automation

`orchestrate-work` → `diagnose-root-cause` → `create-sop` → `automate-routine` → controlled run → `run-retrospective`

Result: a controlled operating process with exceptions, human gates, recovery, and measured improvement.

### Global content delivery

`orchestrate-work` → `write-brief` → `draft-document` → `localize-content` → `write-message` → `review-deliverable`

Result: channel-ready English, Simplified Chinese, and Japanese content that preserves the approved meaning and terminology.

## Capability boundary

AnyWork installs work methods, not credentials or silent integrations. A Skill can use only the tools and access already available to its host agent. External writes, messages, purchases, deletions, and other consequential actions still require explicit authorization and a verifiable target. All 24 atomic Skills and the orchestrator remain `experimental` until the [quality gates](docs/QUALITY.md) and [roadmap](ROADMAP.md) evidence are complete.
