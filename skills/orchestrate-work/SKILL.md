---
name: orchestrate-work
description: "Coordinate a complex outcome by selecting and executing the smallest useful AnyWork Skill chain with checked handoffs, evidence, and permission gates. Use for multi-stage work such as research-to-recommendation, data-to-deck, meeting-to-execution, or process-to-automation; not for one bounded task. 中文：为复杂成果选择并执行最小必要 Skill 链，检查交接、证据与权限；用于调研到建议、数据到汇报、会议到执行、流程到自动化，不用于单一明确任务。日本語：複雑な成果に必要最小限の Skill 連鎖を選択・実行し、引き継ぎ、根拠、権限を検査する。調査から提案、データから資料、会議から実行、手順から自動化に使い、単一の明確なタスクには使わない。"
---

# Orchestrate Work

Turn one complex request into a checked final outcome. Execute the work; do not stop after producing a plan unless the user asked for planning only.

## Workflow

1. Restate the outcome, audience, final deliverable, constraints, acceptance checks, and authorized side effects.
2. Separate known facts, assumptions, missing inputs, and decisions. Ask only for information that blocks a safe or materially correct result.
3. Select the smallest chain of atomic Skills that covers the outcome. Read [references/recipes.md](references/recipes.md) when a standard chain fits or routing is uncertain.
4. Show the chain in one compact line, including the output handed from each stage to the next. Do not turn this into a long planning ceremony.
5. Execute each stage in order. Load and follow the selected Skill before performing that stage; do not merely mention it.
6. Preserve a shared evidence trail: source, timestamp when relevant, assumptions, transformations, decisions, owners, and unresolved risks.
7. Check every handoff before continuing. Reject missing evidence, ambiguous units, invented facts, unresolved goal conflicts, or an output that the next stage cannot consume.
8. Before any external or consequential write, verify the target and obtain explicit authorization. Approval for one stage does not authorize later writes.
9. Review the assembled deliverable against the original acceptance checks. Fix in-scope defects, then return the final artifact, decisions, evidence limits, and next action.

## Routing rules

- Use `clarify-outcome` first only when outcome or acceptance criteria are genuinely unclear.
- Use `gather-sources`, `verify-claims`, and `synthesize-sources` for discovery, claim testing, and synthesis respectively; do not collapse their evidence contracts.
- Use `analyze-data` for observed data and `model-scenarios` for assumption-driven futures.
- Use `compare-options` before `make-decision` when alternatives are not yet evaluated.
- Choose the final production Skill by artifact: `write-brief`, `draft-document`, `build-spreadsheet`, `design-presentation`, or `write-message`.
- Use `prepare-meeting` before a meeting and `capture-decisions-actions` after it.
- Use `create-sop` for a human procedure and `automate-routine` for tool- or code-driven execution.
- Finish with `review-deliverable` when the result will be shared or used for a consequential decision; use `run-retrospective` only after completed work.

## Handoff contract

For every stage, keep a compact record:

| Field | Requirement |
|---|---|
| Input | Named artifact or verified context consumed by the stage |
| Output | Concrete artifact produced for the next stage |
| Evidence | Sources, data, calculations, or observations supporting it |
| Uncertainty | Assumptions, gaps, confidence, and disputed points |
| Gate | Acceptance check and any permission required to continue |

## Boundaries

- Prefer two precise Skills over a large speculative chain.
- Do not claim a Skill was executed unless its workflow was actually followed.
- Do not hide missing tools, access, source material, or model limitations.
- Do not pre-authorize tools or infer permission from the user's desire for an end result.
- If a selected Skill is unavailable, follow its declared artifact contract manually and disclose the fallback.
- Follow the user's language; keep quoted evidence in its source language when translation could change meaning.
