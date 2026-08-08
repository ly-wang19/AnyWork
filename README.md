# AnyWork

**Any work. Any agent. Any machine.**

[简体中文](README.zh-CN.md) · **English** · [日本語](README.ja.md)

[Capabilities](CAPABILITIES.md) · [Benchmark](docs/BENCHMARK.md) · [Roadmap](ROADMAP.md) · [Quality](docs/QUALITY.md) · [Support](docs/SUPPORT.md) · [Contributing](CONTRIBUTING.md)

AnyWork is an Apache-2.0 work-capability layer for AI agents: 24 atomic Skills spanning research, analysis, decisions, writing, data, meetings, operations, and continuous improvement, plus one orchestrator that turns a complex request into the right checked capability chain.

This repository is an **experimental alpha**, not a finished “world’s best” claim. AnyWork makes that ambition testable: every capability, language, host, and machine claim must point to reproducible evidence.

## What AnyWork can do

| Work area | Available outcomes | Skills |
|---|---|---:|
| Define and plan | Clarify an outcome, build a plan, prioritize a backlog | 3 |
| Research and evidence | Gather sources, extract information, verify claims, synthesize evidence | 4 |
| Analyze and decide | Analyze data, diagnose causes, model scenarios, compare options, recommend a decision | 5 |
| Create and communicate | Produce briefs, documents, spreadsheets, presentations, messages, and localized content | 6 |
| Meetings and operations | Prepare meetings, capture decisions, create SOPs, automate routines | 4 |
| Quality and improvement | Review deliverables and run evidence-based retrospectives | 2 |
| End-to-end orchestration | Select, execute, check, and hand off the smallest useful multi-Skill chain | 1 |

Install one of **7 work packs**—Essential, Manager & Leadership, Product & Operations, Research & Consulting, Go-to-Market, People & Recruiting, or Complete—or compose the 24 atomic Skills yourself. See the [complete capability map, all 25 capabilities, and example workflows](CAPABILITIES.md).

## North-star contract

- **Any work:** universal atomic skills plus composable work and role packs.
- **Any agent:** one conservative canonical skill source, with documented host adapters.
- **Any machine:** native binaries, deterministic install/update, recovery, diagnostics, and recoverable uninstall.
- **Three languages:** English, Simplified Chinese, and Japanese share the same structural contract and stable-release thresholds.
- **Trust by default:** provenance, license, permissions, and risk metadata are required.

Current alpha truth: all 24 atomic Skills and the orchestrator are `experimental`; all three language versions are machine-drafted and await signed human review; cross-agent behavioral parity and the full native machine matrix are not yet verified. See the [claims ledger](evidence/claims.json) and [quality gates](docs/QUALITY.md).

Documented default adapters cover Codex, Claude Code, Gemini CLI, GitHub Copilot, and OpenCode. Cursor is experimental and excluded from `all`. The static evaluation contract contains 75 cases and 225 trilingual prompt variants; it does not claim that those prompts have already been executed across models.

## One-command start

macOS or Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/ly-wang19/AnyWork/v0.3.0-alpha.1/install.sh | sh -s -- --setup --agent all --lang en
```

Windows PowerShell:

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/ly-wang19/AnyWork/v0.3.0-alpha.1/install.ps1))) -Setup -Agent all -Language en
```

The version-pinned scripts detect the machine, download the matching native archive, verify its published SHA-256 checksum, install the CLI, and run `anywork setup`. Review [install.sh](install.sh) or [install.ps1](install.ps1) before execution if your policy forbids remote scripts. The scripts never edit shell profiles or silently grant tool permissions.

Then see an actual end-to-end job:

```bash
anywork demo research --lang en
```

Paste the displayed request into Codex, Claude Code, or another installed host. Use `--lang en`, `--lang zh-CN`, or `--lang ja`; the default follows the OS locale. To install one smaller pack instead, run `anywork install essential`. `anywork-cli.py` remains a Python standard-library reference implementation and test oracle.

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
CAPABILITIES*.md  24 atomic Skills, 1 orchestrator, 7 work packs, and workflows
ROADMAP*.md   evidence-gated plan from alpha to stable ecosystem
```

## Roadmap

The `v0.3` last-mile experience adds one-command setup, automatic multi-Skill orchestration, runnable demos, and an exact-prompt [before/after benchmark](docs/BENCHMARK.md). The remaining evidence gates are native-language review, model-executed trilingual results, and signed releases. See the [full evidence-gated roadmap](ROADMAP.md).

## Status

AnyWork is pre-release. Structural catalog checks, local lifecycle tests, three-OS CI, six-target cross-compilation, and the public GitHub repository are verified. Native-speaker review, full model-executed evaluations, native execution on every OS/architecture, and signed releases remain release gates—not implied accomplishments.

## License

Apache-2.0. Third-party catalog entries retain their own licenses and attribution; incompatible sources are linked rather than copied.
