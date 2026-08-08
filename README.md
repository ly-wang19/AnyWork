# AnyWork

**Any work. Any agent. Any machine.**

[简体中文](README.zh-CN.md) · **English** · [日本語](README.ja.md)

AnyWork is an Apache-2.0 work-capability layer for AI agents: 24 atomic skills spanning research, analysis, decisions, writing, data, meetings, operations, and continuous improvement, plus a native cross-platform installer.

This repository is an **experimental alpha**, not a finished “world’s best” claim. AnyWork makes that ambition testable: every capability, language, host, and machine claim must point to reproducible evidence.

## North-star contract

- **Any work:** universal atomic skills plus composable work and role packs.
- **Any agent:** one conservative canonical skill source, with documented host adapters.
- **Any machine:** native binaries, deterministic install/update, recovery, diagnostics, and recoverable uninstall.
- **Three languages:** English, Simplified Chinese, and Japanese share the same structural contract and stable-release thresholds.
- **Trust by default:** provenance, license, permissions, and risk metadata are required.

Current alpha truth: all 24 skills are `experimental`; all three language versions are machine-drafted and await signed human review; cross-agent behavioral parity and the full native machine matrix are not yet verified. See the [claims ledger](evidence/claims.json) and [quality gates](docs/QUALITY.md).

Documented default adapters cover Codex, Claude Code, Gemini CLI, GitHub Copilot, and OpenCode. Cursor is experimental and excluded from `all`. The static evaluation contract contains 72 cases and 216 trilingual prompt variants; it does not claim that those prompts have already been executed across models.

## Quick start

Build the dependency-free native CLI from source:

```bash
go build -trimpath -o anywork .
./anywork list
./anywork install essential --scope project --dry-run
./anywork install essential --scope project
./anywork doctor
```

Use `--lang en`, `--lang zh-CN`, or `--lang ja`; the default follows the OS locale. On Windows, build `anywork.exe`. Tagged releases are configured to produce checksum-protected archives for macOS, Linux, and Windows on amd64 and arm64. `anywork-cli.py` remains a Python standard-library reference implementation and test oracle.

## Repository map

```text
skills/       canonical, task-focused skills
registry/     packs plus source, language-review, capability, risk, and version metadata
adapters/     agent-specific installation rules
evals/        multilingual behavioral and quality evaluations
evidence/     machine-readable claim status and supporting artifacts
*.go          dependency-free native AnyWork CLI
src/          Python reference implementation
schemas/      machine-readable contracts
scripts/      reproducible release builds
tests/        registry, evidence, i18n, and installer regression tests
```

## Status

AnyWork is pre-release. Structural catalog checks, local lifecycle tests, and six-target cross-compilation pass. Native-speaker review, full model-executed evaluations, native execution on every OS/architecture, a public GitHub remote, and signed releases remain release gates—not implied accomplishments.

## License

Apache-2.0. Third-party catalog entries retain their own licenses and attribution; incompatible sources are linked rather than copied.
