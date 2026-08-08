# AnyWork roadmap

[简体中文](ROADMAP.zh-CN.md) · **English** · [日本語](ROADMAP.ja.md)

AnyWork is aiming for the most trusted portable work-capability layer for AI agents. This roadmap is evidence-gated: the order is intentional, but versions are not calendar promises. A milestone is complete only when its exit evidence is public and reproducible.

## Status legend

- ✅ **Complete:** implemented and backed by repository evidence.
- 🚧 **In progress:** active work with incomplete exit evidence.
- ⏳ **Planned:** sequenced, not yet started or not yet evidenced.

## Milestones

| Milestone | Outcome | Exit criteria | Status |
|---|---|---|---|
| `v0.2` Portable alpha | A safe, trilingual core that can move between agents and machines | 24 atomic Skills, 7 packs, 5 default host adapters, recoverable native installer, three-OS CI, six-target builds, public claims ledger | ✅ Complete |
| `v0.3` Evidence alpha | Make the layer usable in one step, then replace structural coverage with measured behavioral evidence | One-command setup, end-to-end orchestration and demos, exact-prompt A/B runner; signed native-speaker review; model-executed trilingual evaluations; native artifact runs on all six targets; published raw results | 🚧 In progress |
| `v0.5` Work beta | Cover complete work systems, not isolated prompts | Role packs, office artifacts, research, communication, meetings, operations, people, sales, customer, finance, and compliance workflows; connector boundaries and permission tests | ⏳ Planned |
| `v0.8` Ecosystem beta | Make installation, contribution, and trust scalable | Homebrew/WinGet/Scoop or equivalent distribution, signed registry metadata, community-pack review path, compatibility dashboard, migration guarantees | ⏳ Planned |
| `v1.0` Stable | A dependable default for new machines and professional work | Stable Skills meet every quality threshold; signed releases, SBOM and provenance; supported upgrade policy; no critical safety regressions; published support SLA | ⏳ Planned |
| Beyond `v1.0` | An open work-capability ecosystem | Verified third-party packs, organization policies, private registries, benchmark federation, and sustainable governance | ⏳ Planned |

## Current priorities

1. Complete qualified human review for English, Simplified Chinese, and Japanese against the same Skill revisions.
2. Execute the 225 prompt variants with the exact-prompt A/B runner on the supported agent matrix and publish raw, reproducible results without answer leakage.
3. Run release artifacts natively on macOS, Linux, and Windows across amd64 and arm64, including installation, update, rollback, recovery, and uninstall.
4. Publish the first signed pre-release with checksums, SBOM, provenance attestation, and verified archive contents.
5. Validate the one-command installers on clean machines and add managed upgrade channels.
6. Expand the orchestrator from reusable recipes into end-to-end role and work-system packs while keeping external writes permission-gated.

## Stable-release scorecard

`v1.0` cannot be declared from feature count alone. It requires:

- 100% schema, packaging, critical-safety, and deterministic-fixture gates.
- At least 98% same-language adherence.
- At least 90% trigger recall and at most 5% false-positive rate.
- At least 85/100 average task quality in each language.
- At most a five-point quality gap across the three languages.
- Verified lifecycle behavior for every supported host and native platform class.
- Zero unresolved critical security, provenance, licensing, or data-loss issues.

The authoritative thresholds live in the [quality standard](docs/QUALITY.md); current evidence lives in the [claims ledger](evidence/claims.json).

## Non-negotiables

- No “world's best,” stable, native-language, or cross-agent parity claim without matching public evidence.
- No hidden telemetry, secret collection, install-time execution, or silent tool pre-approval.
- No destructive or external write without explicit authorization and a verifiable target.
- No copied capability without compatible licensing, immutable provenance, and attribution.

## Contributing to the roadmap

Proposals should name the user outcome, affected hosts and languages, permissions or side effects, measurable exit criteria, and the evidence artifact that would prove completion. See [CONTRIBUTING.md](CONTRIBUTING.md).
