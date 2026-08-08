# 智能体与机器支持

**简体中文** · [English](SUPPORT.md) · [日本語](SUPPORT.ja.md)

“支持”包含三个互相独立的层次：官方路径兼容、安装生命周期验证、模型行为等价。AnyWork 必须分别报告，不能混为一个结论。

| 智能体 | 适配器 | 目标族 | 默认安装 | 行为等价 |
|---|---|---|---|---|
| Codex | 稳定 | `.agents/skills` | 是 | 未验证 |
| Claude Code | 稳定 | `.claude/skills` | 是 | 未验证 |
| Gemini CLI | 稳定 | `.agents/skills` | 是 | 未验证 |
| GitHub Copilot | 稳定 | `.agents/skills` | 是 | 未验证 |
| OpenCode | 稳定 | `.agents/skills` | 是 | 未验证 |
| Cursor | 实验 | `.cursor/skills` | 否 | 未验证 |

四个使用 `.agents/skills` 的宿主共享一份物理安装。状态文件按解析后的目标记录消费者集合；卸载一个消费者时，不得删除其他消费者仍在使用的文件。Cursor 在真实宿主上完成路径语义、IDE 发现和 CLI 发现验证前，只能显式选择，不能进入默认安装。

标准 `SKILL.md` frontmatter 只使用 Agent Skills 的可移植字段：`name`、`description`、`license`、`compatibility` 和字符串到字符串的 `metadata`。AnyWork 绝不把风险声明转换成 `allowed-tools` 或任何宿主免审批机制。

原生压缩包面向 macOS、Linux、Windows 的 amd64 与 arm64。六目标交叉编译已经验证；六类原生 runner 上实跑发布产物仍是发布门槛。权威状态见 [evidence/claims.json](../evidence/claims.json)，适配器细节见 [registry/support-matrix.json](../registry/support-matrix.json)。
