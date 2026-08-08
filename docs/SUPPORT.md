# Host and machine support

[简体中文](SUPPORT.zh-CN.md) · **English** · [日本語](SUPPORT.ja.md)

“Supported” has three independent levels: documented path compatibility, installer lifecycle verification, and model-behavior parity. AnyWork reports them separately.

| Host | Adapter | Target family | Default install | Behavioral parity |
|---|---|---|---|---|
| Codex | stable | `.agents/skills` | yes | not verified |
| Claude Code | stable | `.claude/skills` | yes | not verified |
| Gemini CLI | stable | `.agents/skills` | yes | not verified |
| GitHub Copilot | stable | `.agents/skills` | yes | not verified |
| OpenCode | stable | `.agents/skills` | yes | not verified |
| Cursor | experimental | `.cursor/skills` | no | not verified |

The four `.agents/skills` hosts share one physical installation. State tracks a set of consumers for each resolved target; uninstalling one consumer must not remove files still used by another. Cursor remains opt-in until its exact path semantics and both IDE/CLI discovery are verified on a real host.

Canonical `SKILL.md` frontmatter uses only the portable Agent Skills fields `name`, `description`, `license`, `compatibility`, and string-to-string `metadata`. AnyWork never converts risk declarations into `allowed-tools` or another host pre-approval mechanism.

Native archives are built for macOS, Linux, and Windows on amd64 and arm64. Cross-compilation is verified; execution of the release artifact on all six native runner classes remains a release gate. The authoritative status is [evidence/claims.json](../evidence/claims.json), and descriptor details live in [registry/support-matrix.json](../registry/support-matrix.json).
