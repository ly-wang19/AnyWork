---
name: analyze-data
description: "Inspect existing data and produce reproducible findings after quality checks. Use for datasets, metrics, experiments, and operational analysis; do not use to build a forecasting scenario or a spreadsheet deliverable. 中文：检查已有数据质量并产出可复算结论；适用于数据集、指标、实验和运营分析，不用于建立预测情景或交付电子表格。日本語：既存データの品質を確認し、再現可能な知見を出す。データセット、指標、実験、業務分析に使い、予測シナリオやスプレッドシート成果物の作成には使わない。"
---

# Analyze Data

Produce an evidence-aware, directly usable work result. Keep facts, assumptions, and recommendations distinct.

## Workflow

1. Confirm the business question, grain, units, time zone, and metric definitions.
2. Profile missing values, duplicates, outliers, invalid types, and coverage gaps.
3. Document cleaning and exclusion decisions before calculating results.
4. Use reproducible formulas or code and preserve intermediate checks.
5. Separate observations, statistical evidence, interpretation, and recommendation.
6. Test sensitivity to material assumptions and disclose limitations.

## Deliver

- Question and metric definitions
- Data-quality report
- Method and calculations
- Findings and visual evidence
- Limitations and next tests

## Quality gates

- Do not hide inconvenient rows or silently change metric definitions.
- Do not imply causation from association without a valid design.

## Language

- Follow the user's latest explicit output-language instruction.
- Otherwise answer in the dominant language of the request: English, Simplified Chinese, or Japanese.
- Preserve quoted source text, code, identifiers, product names, and citations in their original language unless translation is requested.
- Match the audience's register; distinguish formal and conversational Chinese, business and plain Japanese, and the requested English variety.
