from __future__ import annotations

import locale
import os
from typing import Any


SUPPORTED_LANGUAGES = ("en", "zh-CN", "ja")

MESSAGES: dict[str, dict[str, str]] = {
    "en": {
        "catalog_title": "Available AnyWork packs",
        "pack_line": "{name}: {description} ({count} skills)",
        "plan_install": "PLAN install {skill} -> {target}",
        "plan_replace": "PLAN replace {skill} -> {target}",
        "plan_remove": "PLAN uninstall {skill} from {target}",
        "installed": "Installed {skill} for {agent}: {target}",
        "updated": "Updated {skill} for {agent}: {target}",
        "unchanged": "Already current: {skill} for {agent}",
        "removed": "Uninstalled {skill} for {agent}; backup: {backup}",
        "conflict": "Refusing to overwrite modified or unmanaged skill: {target}. Use --force after review.",
        "not_owned": "Refusing to uninstall an unmanaged skill: {target}",
        "pack_unknown": "Unknown pack: {pack}",
        "agent_unknown": "Unsupported agent: {agent}",
        "doctor_ok": "AnyWork doctor passed: {skills} skills, {packs} packs, {agents} agents.",
        "doctor_fail": "AnyWork doctor found {count} problem(s):",
        "validation_error": "- {message}",
        "nothing_installed": "No AnyWork-managed installation found for {skill} on {agent}.",
        "backup_note": "Existing content backed up to {backup}",
    },
    "zh-CN": {
        "catalog_title": "可用的 AnyWork 工作包",
        "pack_line": "{name}：{description}（{count} 个技能）",
        "plan_install": "计划安装 {skill} -> {target}",
        "plan_replace": "计划更新 {skill} -> {target}",
        "plan_remove": "计划卸载 {skill}：{target}",
        "installed": "已为 {agent} 安装 {skill}：{target}",
        "updated": "已为 {agent} 更新 {skill}：{target}",
        "unchanged": "已是最新版本：{agent} / {skill}",
        "removed": "已为 {agent} 卸载 {skill}；备份：{backup}",
        "conflict": "拒绝覆盖已修改或非 AnyWork 管理的技能：{target}。请审查后使用 --force。",
        "not_owned": "拒绝卸载非 AnyWork 管理的技能：{target}",
        "pack_unknown": "未知工作包：{pack}",
        "agent_unknown": "不支持的智能体：{agent}",
        "doctor_ok": "AnyWork 诊断通过：{skills} 个技能、{packs} 个工作包、{agents} 个智能体。",
        "doctor_fail": "AnyWork 诊断发现 {count} 个问题：",
        "validation_error": "- {message}",
        "nothing_installed": "未找到由 AnyWork 管理的安装：{agent} / {skill}。",
        "backup_note": "原内容已备份到 {backup}",
    },
    "ja": {
        "catalog_title": "利用可能な AnyWork パック",
        "pack_line": "{name}：{description}（{count} スキル）",
        "plan_install": "インストール予定 {skill} -> {target}",
        "plan_replace": "更新予定 {skill} -> {target}",
        "plan_remove": "アンインストール予定 {skill}：{target}",
        "installed": "{agent} に {skill} をインストールしました：{target}",
        "updated": "{agent} の {skill} を更新しました：{target}",
        "unchanged": "最新版です：{agent} / {skill}",
        "removed": "{agent} から {skill} を削除しました。バックアップ：{backup}",
        "conflict": "変更済み、または AnyWork 管理外のスキルは上書きしません：{target}。確認後に --force を使用してください。",
        "not_owned": "AnyWork 管理外のスキルはアンインストールしません：{target}",
        "pack_unknown": "不明なパック：{pack}",
        "agent_unknown": "未対応のエージェント：{agent}",
        "doctor_ok": "AnyWork の診断に合格しました：{skills} スキル、{packs} パック、{agents} エージェント。",
        "doctor_fail": "AnyWork の診断で {count} 件の問題が見つかりました：",
        "validation_error": "- {message}",
        "nothing_installed": "AnyWork が管理するインストールはありません：{agent} / {skill}。",
        "backup_note": "既存内容を {backup} にバックアップしました",
    },
}

CLI_HELP: dict[str, dict[str, str]] = {
    "en": {
        "description": "Install and manage trusted work capabilities across AI agents.",
        "help": "show this help message and exit",
        "list": "list available packs",
        "doctor": "validate the catalog, skills, evaluations, and safety rules",
        "install": "install a pack",
        "update": "update a managed pack",
        "uninstall": "uninstall a managed pack into a recoverable backup",
        "pack": "pack id",
        "agent": "codex, claude-code, or all; repeat to select multiple agents",
        "scope": "install for the current user or one project",
        "project_dir": "project directory used with project scope",
        "dry_run": "show the plan without writing files",
        "force": "replace drifted or unmanaged content after review",
    },
    "zh-CN": {
        "description": "跨 AI 智能体安装和管理可信工作能力。",
        "help": "显示帮助并退出",
        "list": "列出可用工作包",
        "doctor": "校验目录、技能、评测和安全规则",
        "install": "安装工作包",
        "update": "更新由 AnyWork 管理的工作包",
        "uninstall": "卸载工作包并保留可恢复备份",
        "pack": "工作包 ID",
        "agent": "codex、claude-code 或 all；可重复指定",
        "scope": "安装到当前用户或一个项目",
        "project_dir": "项目级安装所使用的项目目录",
        "dry_run": "只显示计划，不写入文件",
        "force": "审查后替换已漂移或非托管内容",
    },
    "ja": {
        "description": "信頼できる仕事能力を複数の AI Agent にインストールして管理します。",
        "help": "ヘルプを表示して終了",
        "list": "利用可能なパックを表示",
        "doctor": "カタログ、Skill、評価、安全規則を検証",
        "install": "パックをインストール",
        "update": "AnyWork 管理パックを更新",
        "uninstall": "復元可能なバックアップを残してパックを削除",
        "pack": "パック ID",
        "agent": "codex、claude-code、all のいずれか。複数回指定可能",
        "scope": "現在のユーザー、または一つのプロジェクトへ導入",
        "project_dir": "プロジェクトスコープで使うディレクトリ",
        "dry_run": "ファイルを書き込まず計画のみ表示",
        "force": "確認後、変更済みまたは管理外の内容を置換",
    },
}


def normalize_language(value: str | None) -> str:
    if value:
        normalized = value.replace("_", "-")
        lowered = normalized.lower()
        if lowered.startswith("zh"):
            return "zh-CN"
        if lowered.startswith("ja"):
            return "ja"
        return "en"
    configured = os.environ.get("ANYWORK_LOCALE")
    if configured:
        return normalize_language(configured)
    detected = locale.getlocale()[0] or "en"
    return normalize_language(detected)


def message(language: str, key: str, **values: Any) -> str:
    template = MESSAGES.get(language, MESSAGES["en"]).get(key, MESSAGES["en"].get(key, key))
    return template.format(**values)


def localized(value: Any, language: str) -> str:
    if isinstance(value, str):
        return value
    if not isinstance(value, dict):
        return str(value)
    return str(value.get(language) or value.get("en") or next(iter(value.values()), ""))


def help_text(language: str, key: str) -> str:
    return CLI_HELP.get(language, CLI_HELP["en"]).get(key, CLI_HELP["en"].get(key, key))
