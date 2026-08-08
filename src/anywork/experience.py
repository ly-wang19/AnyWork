from __future__ import annotations

from dataclasses import dataclass

from .i18n import localized, message


@dataclass(frozen=True)
class Demo:
    id: str
    name: dict[str, str]
    prompt: dict[str, str]
    chain: tuple[str, ...]
    outcome: dict[str, str]


DEMOS = (
    Demo(
        id="research",
        name={"en": "Research to recommendation", "zh-CN": "从调研到建议", "ja": "調査から提案まで"},
        prompt={
            "en": "Use $orchestrate-work to evaluate three employee knowledge-base options end to end. Gather current primary sources, verify material claims, compare total cost, security, migration effort, and adoption risk, recommend one option, and deliver a decision-ready brief with citations and explicit unknowns. Do not contact vendors or purchase anything.",
            "zh-CN": "使用 $orchestrate-work 端到端评估三个员工知识库方案。查找最新一手来源，核验重要主张，比较总成本、安全性、迁移工作量和采用风险，推荐一个方案，并交付包含引用和明确未知项的决策简报。不要联系供应商或进行购买。",
            "ja": "$orchestrate-work を使い、従業員向けナレッジベース3案を最後まで評価して。最新の一次情報を集め、重要な主張を検証し、総費用、セキュリティ、移行工数、定着リスクを比較して1案を提案し、引用と不明点を明記した意思決定用ブリーフを作成して。ベンダーへの連絡や購入は行わないで。",
        },
        chain=(
            "orchestrate-work",
            "gather-sources",
            "verify-claims",
            "compare-options",
            "make-decision",
            "write-brief",
            "review-deliverable",
        ),
        outcome={
            "en": "An auditable vendor recommendation and checked decision brief.",
            "zh-CN": "可审计的供应商建议和经过检查的决策简报。",
            "ja": "監査可能なベンダー提案と検証済み意思決定ブリーフ。",
        },
    ),
    Demo(
        id="meeting",
        name={"en": "Meeting to accountable execution", "zh-CN": "从会议到责任落实", "ja": "会議から責任ある実行まで"},
        prompt={
            "en": "Use $orchestrate-work to turn this product planning session into accountable execution. Prepare the decision agenda and pre-read, then after I provide the notes, extract only supported decisions and actions, assign only stated owners and dates, build the delivery plan, and flag unresolved dependencies.",
            "zh-CN": "使用 $orchestrate-work 把这次产品规划会转成责任到人的执行。先准备决策议程和预读；我提供会议记录后，只提取有依据的决策与行动，只记录明确说明的负责人和日期，形成交付计划，并标出未解决依赖。",
            "ja": "$orchestrate-work を使い、この製品計画会議を責任ある実行につなげて。まず意思決定議題と事前資料を準備し、私がメモを渡した後、根拠のある決定とアクションだけを抽出し、明記された担当者と期限だけを記録して、実行計画と未解決の依存関係を整理して。",
        },
        chain=(
            "orchestrate-work",
            "prepare-meeting",
            "capture-decisions-actions",
            "plan-work",
            "prioritize-work",
            "review-deliverable",
        ),
        outcome={
            "en": "A decision-focused meeting and a traceable owned execution plan.",
            "zh-CN": "聚焦决策的会议和可追溯、责任明确的执行计划。",
            "ja": "意思決定に集中した会議と、追跡可能で担当が明確な実行計画。",
        },
    ),
    Demo(
        id="automation",
        name={"en": "Recurring operation to safe automation", "zh-CN": "从重复运营到安全自动化", "ja": "反復運用から安全な自動化まで"},
        prompt={
            "en": "Use $orchestrate-work to eliminate recurring weekly-report failures. Diagnose the root cause from the logs and examples, write a controlled SOP with exceptions and escalation, implement an idempotent automation with dry-run, least privilege, logging, and rollback, test it locally, and return the evidence. Do not deploy or schedule it without approval.",
            "zh-CN": "使用 $orchestrate-work 消除每周报告反复失败的问题。根据日志和样例诊断根因，编写包含异常与升级路径的受控 SOP，实现带预演、最小权限、日志和回滚的幂等自动化，在本地测试并返回证据。未经批准不得部署或设置定时任务。",
            "ja": "$orchestrate-work を使い、週次レポートの反復障害を解消して。ログと例から根本原因を診断し、例外とエスカレーションを含む統制済みSOPを作成して、ドライラン、最小権限、ログ、ロールバックを備えた冪等な自動化を実装し、ローカルで検証して証拠を返して。承認なしに導入やスケジュール設定は行わないで。",
        },
        chain=(
            "orchestrate-work",
            "diagnose-root-cause",
            "create-sop",
            "automate-routine",
            "review-deliverable",
            "run-retrospective",
        ),
        outcome={
            "en": "A tested, recoverable automation and its operating controls.",
            "zh-CN": "经过测试、可恢复的自动化及其运营控制。",
            "ja": "検証済みで復旧可能な自動化と運用統制。",
        },
    ),
)


def render_demo(language: str, demo_id: str | None = None) -> str:
    if not demo_id:
        lines = [message(language, "demo_title")]
        lines.extend(f"- {demo.id}: {localized(demo.name, language)}" for demo in DEMOS)
        lines.append(message(language, "demo_choose"))
        return "\n".join(lines)
    selected = next((demo for demo in DEMOS if demo.id == demo_id), None)
    if selected is None:
        raise KeyError(demo_id)
    return "\n\n".join(
        (
            localized(selected.name, language),
            f"{message(language, 'demo_prompt')}\n{localized(selected.prompt, language)}",
            f"{message(language, 'demo_chain')}\n{' → '.join(selected.chain)}",
            f"{message(language, 'demo_outcome')}\n{localized(selected.outcome, language)}",
        )
    )
