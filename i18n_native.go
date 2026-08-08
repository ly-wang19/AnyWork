package main

import (
	"fmt"
	"os"
	"strings"
)

var nativeMessages = map[string]map[string]string{
	"en": {
		"catalog_title":       "Available AnyWork packs",
		"quality_notice":      "Alpha quality: all 24 skills are experimental; EN/ZH-CN/JA content is machine-drafted and awaits signed human review.",
		"doctor_ok":           "AnyWork doctor passed: %d skills, %d packs, %d agents.",
		"doctor_fail":         "AnyWork doctor found %d problem(s):",
		"plan_install":        "PLAN install %s -> %s",
		"plan_replace":        "PLAN replace %s -> %s",
		"plan_remove":         "PLAN uninstall %s from %s",
		"installed":           "Installed %s for %s: %s",
		"updated":             "Updated %s for %s: %s",
		"unchanged":           "Already current: %s for %s",
		"removed":             "Uninstalled %s for %s; backup: %s",
		"backup":              "Existing content backed up to %s",
		"conflict":            "Refusing to overwrite modified or unmanaged skill: %s. Use --force after review.",
		"not_owned":           "Refusing to uninstall modified or unmanaged skill: %s",
		"unknown_pack":        "Unknown pack: %s",
		"lock_busy":           "Another AnyWork mutation is active: %s",
		"recover_none":        "No unfinished AnyWork transaction was found.",
		"recover_committed":   "Recovered committed transaction %s.",
		"recover_rolled_back": "Rolled back unfinished transaction %s; displaced content was preserved in the adjacent recovery area.",
		"usage":               "Usage: anywork <list|doctor|install|update|uninstall|recover> [options]",
	},
	"zh-CN": {
		"catalog_title":       "可用的 AnyWork 工作包",
		"quality_notice":      "Alpha 质量状态：24 个技能均为实验版；中、英、日内容由机器起草，尚待人工签署审校。",
		"doctor_ok":           "AnyWork 诊断通过：%d 个技能、%d 个工作包、%d 个智能体。",
		"doctor_fail":         "AnyWork 诊断发现 %d 个问题：",
		"plan_install":        "计划安装 %s -> %s",
		"plan_replace":        "计划更新 %s -> %s",
		"plan_remove":         "计划卸载 %s：%s",
		"installed":           "已安装 %s（%s）：%s",
		"updated":             "已更新 %s（%s）：%s",
		"unchanged":           "已是最新版本：%s / %s",
		"removed":             "已卸载 %s（%s）；备份：%s",
		"backup":              "原内容已备份到 %s",
		"conflict":            "拒绝覆盖已修改或非 AnyWork 管理的技能：%s。请审查后使用 --force。",
		"not_owned":           "拒绝卸载已修改或非 AnyWork 管理的技能：%s",
		"unknown_pack":        "未知工作包：%s",
		"lock_busy":           "另一个 AnyWork 写入操作正在进行：%s",
		"recover_none":        "未发现未完成的 AnyWork 事务。",
		"recover_committed":   "已恢复完成提交的事务 %s。",
		"recover_rolled_back": "已回滚未完成的事务 %s；被替换的内容已保留在相邻恢复目录。",
		"usage":               "用法：anywork <list|doctor|install|update|uninstall|recover> [选项]",
	},
	"ja": {
		"catalog_title":       "利用可能な AnyWork パック",
		"quality_notice":      "Alpha 品質：24 スキルはすべて実験版です。英語・中国語・日本語の内容は機械生成の初稿で、署名付きの人手レビューは未完了です。",
		"doctor_ok":           "AnyWork の診断は正常です：%d スキル、%d パック、%d エージェント。",
		"doctor_fail":         "AnyWork の診断で %d 件の問題が見つかりました：",
		"plan_install":        "インストール予定 %s -> %s",
		"plan_replace":        "更新予定 %s -> %s",
		"plan_remove":         "アンインストール予定 %s：%s",
		"installed":           "%s をインストールしました（%s）：%s",
		"updated":             "%s を更新しました（%s）：%s",
		"unchanged":           "最新版です：%s / %s",
		"removed":             "%s をアンインストールしました（%s）。バックアップ：%s",
		"backup":              "既存内容を %s にバックアップしました",
		"conflict":            "変更済み、または AnyWork 管理外のスキルは上書きしません：%s。確認後に --force を使用してください。",
		"not_owned":           "変更済み、または AnyWork 管理外のスキルはアンインストールしません：%s",
		"unknown_pack":        "不明なパック：%s",
		"lock_busy":           "別の AnyWork 更新処理が実行中です：%s",
		"recover_none":        "未完了の AnyWork トランザクションはありません。",
		"recover_committed":   "コミット済みトランザクション %s を復旧しました。",
		"recover_rolled_back": "未完了のトランザクション %s をロールバックしました。退避した内容は隣接する復旧領域に保持されています。",
		"usage":               "使い方：anywork <list|doctor|install|update|uninstall|recover> [オプション]",
	},
}

func normalizeLanguage(value string) string {
	if value == "" {
		value = os.Getenv("ANYWORK_LOCALE")
	}
	if value == "" {
		value = os.Getenv("LC_ALL")
	}
	if value == "" {
		value = os.Getenv("LANG")
	}
	lower := strings.ToLower(strings.ReplaceAll(value, "_", "-"))
	switch {
	case strings.HasPrefix(lower, "zh"):
		return "zh-CN"
	case strings.HasPrefix(lower, "ja"):
		return "ja"
	default:
		return "en"
	}
}

func msg(language, key string, values ...any) string {
	template := nativeMessages[language][key]
	if template == "" {
		template = nativeMessages["en"][key]
	}
	if template == "" {
		template = key
	}
	return fmt.Sprintf(template, values...)
}

func localize(value localizedText, language string) string {
	if value[language] != "" {
		return value[language]
	}
	return value["en"]
}

func extractLanguage(arguments []string) (string, []string, error) {
	selected := ""
	cleaned := []string{}
	for index := 0; index < len(arguments); index++ {
		current := arguments[index]
		if current == "--lang" {
			if index+1 >= len(arguments) {
				return "", nil, fmt.Errorf("--lang requires a value")
			}
			selected = arguments[index+1]
			index++
			continue
		}
		if strings.HasPrefix(current, "--lang=") {
			selected = strings.TrimPrefix(current, "--lang=")
			continue
		}
		cleaned = append(cleaned, current)
	}
	if selected != "" {
		prefix := strings.Split(strings.ToLower(strings.ReplaceAll(selected, "_", "-")), "-")[0]
		if prefix != "en" && prefix != "zh" && prefix != "ja" {
			return "", nil, fmt.Errorf("unsupported language: %s", selected)
		}
	}
	return normalizeLanguage(selected), cleaned, nil
}
