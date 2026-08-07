---
name: analyze-data
description: "Inspect data quality and produce reproducible analysis, calculations, and bounded insights. Use for tables, metrics, experiments, or operational data. 中文：先检查数据质量，再做可复算分析和洞察。日本語：データ品質を確認してから、再現可能な分析と洞察を作る。"
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
