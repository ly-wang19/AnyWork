---
name: automate-routine
description: "Design safe and idempotent automation for recurring work with permissions, dry-run, logs, recovery, and human gates. 中文：安全自动化重复工作，包含权限、预演、日志、恢复和人工闸门。日本語：反復業務を権限、ドライラン、ログ、復旧、人の承認付きで安全に自動化する。"
---

# Automate Routine

Produce an evidence-aware, directly usable work result. Keep facts, assumptions, and recommendations distinct.

## Workflow

1. Map the manual process, trigger, inputs, outputs, systems, owners, and failure cost.
2. Separate deterministic steps from judgment and consequential external actions.
3. Define least-privilege capabilities, secret handling, rate limits, and approvals.
4. Make operations idempotent and add dry-run, logging, retry limits, and checkpoints.
5. Design rollback or compensating actions for every material mutation.
6. Test normal, edge, partial-failure, and replay scenarios in a safe environment.

## Deliver

- Automation boundary
- Trigger and data flow
- Permissions and approvals
- Failure and recovery design
- Observability and runbook
- Test evidence

## Quality gates

- Never add silent external writes or install-time code execution.
- Require explicit approval for newly expanded capabilities.

## Language

- Follow the user's latest explicit output-language instruction.
- Otherwise answer in the dominant language of the request: English, Simplified Chinese, or Japanese.
- Preserve quoted source text, code, identifiers, product names, and citations in their original language unless translation is requested.
- Match the audience's register; distinguish formal and conversational Chinese, business and plain Japanese, and the requested English variety.
