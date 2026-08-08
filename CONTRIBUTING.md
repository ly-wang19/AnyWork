# Contributing to AnyWork

[简体中文](CONTRIBUTING.zh-CN.md) · **English** · [日本語](CONTRIBUTING.ja.md)

Contribute task-focused, attributable, safe, and testable work capabilities. Do not submit prompt dumps or copy publicly visible material without a compatible license.

## Before a pull request

1. Choose one atomic work outcome and a lowercase kebab-case Skill name.
2. Keep workflow logic in one `SKILL.md`; cover English, Simplified Chinese, and Japanese in its trigger description and evaluations.
3. Add source, license, version, capabilities, risk, side effects, and honest language-review status to `registry/catalog.json`.
4. For adapted work, record the immutable upstream revision, authors, license, attribution, and modifications in `THIRD_PARTY.yml`.
5. Add multilingual evaluation cases and deterministic tests where applicable.
6. Run:

   ```bash
   go test ./...
   go run . doctor --lang en
   python3 anywork-cli.py doctor
   PYTHONPATH=src python3 -m unittest discover -s tests -v
   ```

7. Sign every commit with `Signed-off-by: Name <email>` to certify the Developer Certificate of Origin 1.1 in [DCO](DCO).

Automated translation may be used as a draft, but stable releases require qualified human review in all three languages. AI-assisted contributions remain the contributor's responsibility for rights, facts, safety, and testing.

## Third-party boundary

- `original`: contributed under Apache-2.0 and DCO.
- `adapted`: copied or modified only under a compatible license with full provenance.
- `curated`: metadata and upstream link only; no copied content.

Unlicensed, source-available, non-commercial, no-derivatives, or otherwise incompatible material must not be copied into the core distribution.
