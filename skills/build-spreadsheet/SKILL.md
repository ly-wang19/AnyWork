---
name: build-spreadsheet
description: "Create a recalculable, auditable spreadsheet with explicit inputs, formulas, outputs, and checks. Use when the deliverable is a workbook, budget, model, or tracker; do not use only to interpret an existing dataset or reason about scenarios without building the file. 中文：创建输入、公式、输出和校验清晰的可复算电子表格；适用于交付工作簿、预算、模型或跟踪表，不用于仅解读已有数据或只做情景推演。日本語：入力、数式、出力、検証が明確な再計算可能スプレッドシートを作る。成果物がワークブック、予算、モデル、管理表の場合に使い、既存データの解釈だけ、またはファイルを作らないシナリオ検討には使わない。"
---

# Build Spreadsheet

Produce an evidence-aware, directly usable work result. Keep facts, assumptions, and recommendations distinct.

## Workflow

1. Define users, decisions, inputs, outputs, units, and refresh process.
2. Separate inputs, calculations, outputs, and reference data.
3. Use formulas rather than hard-coded derived values.
4. Add validation, error handling, totals, and reconciliation checks.
5. Apply readable formats, frozen labels, and accessible color use.
6. Change representative inputs and verify downstream recalculation.

## Deliver

- Workbook structure
- Documented inputs and assumptions
- Auditable formulas
- Decision-ready outputs
- Validation and reconciliation checks

## Quality gates

- Never hide critical logic in unexplained constants.
- Preserve source data and do not overwrite it during cleaning.

## Language

- Follow the user's latest explicit output-language instruction.
- Otherwise answer in the dominant language of the request: English, Simplified Chinese, or Japanese.
- Preserve quoted source text, code, identifiers, product names, and citations in their original language unless translation is requested.
- Match the audience's register; distinguish formal and conversational Chinese, business and plain Japanese, and the requested English variety.
