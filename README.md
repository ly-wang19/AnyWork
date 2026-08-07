# AnyWork

**Any work. Any agent. Any machine.**

[简体中文](README.zh-CN.md) · **English** · [日本語](README.ja.md)

AnyWork is an open-source work capability layer for AI agents. It packages tested workflows, references, scripts, templates, and evaluations so people can restore a high-quality AI workspace on a new machine in minutes.

AnyWork is not a prompt dump. Official skills must be task-focused, multilingual, evaluated, attributable, and safe to install.

## Product contract

- **Any work:** universal work skills plus composable role and industry packs.
- **Any agent:** one canonical skill source with adapters for Codex and Claude Code first.
- **Any machine:** deterministic install, update, doctor, rollback, and uninstall flows.
- **Three languages:** English, Simplified Chinese, and Japanese are first-class across docs, triggers, examples, CLI messages, and evaluations.
- **Trust by default:** provenance, license, permissions, and risk metadata are required.

## Quick start

```bash
python3 anywork-cli.py list
python3 anywork-cli.py install essential --agent all --dry-run
python3 anywork-cli.py install essential --agent all
python3 anywork-cli.py doctor
```

Use `--lang en`, `--lang zh-CN`, or `--lang ja` to select CLI output. The default follows the operating-system locale.
On Windows, use `py anywork-cli.py ...`. The repository CLI has no third-party runtime dependencies.

## Repository map

```text
skills/       canonical, task-focused skills
packs/        composable work and role packs
registry/     source, license, language, risk, and version metadata
adapters/     agent-specific installation rules
evals/        multilingual behavioral and quality evaluations
src/          cross-platform AnyWork CLI
schemas/      machine-readable contracts
tests/        installer and registry regression tests
```

## Status

AnyWork is in its foundation phase. Codex and Claude Code are the first supported adapters. The public quality bar and evaluation suite are part of the product, not post-launch documentation.

## License

Apache-2.0. Third-party catalog entries retain their own licenses and attribution; incompatible sources are linked rather than copied.
