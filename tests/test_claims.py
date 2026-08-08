from __future__ import annotations

import json
import unittest

from anywork.catalog import repository_root


class ClaimsLedgerTests(unittest.TestCase):
    def test_claims_are_unique_and_evidence_paths_exist(self) -> None:
        root = repository_root()
        payload = json.loads((root / "evidence" / "claims.json").read_text(encoding="utf-8"))
        claims = payload["claims"]
        self.assertEqual(len({claim["id"] for claim in claims}), len(claims))
        for claim in claims:
            self.assertIn(claim["status"], {"verified", "partial", "not-verified"})
            if claim["status"] == "verified":
                self.assertTrue(claim["evidence"], claim["id"])
            for relative in claim["evidence"]:
                self.assertTrue((root / relative).exists(), f"missing evidence: {relative}")

    def test_unverified_quality_claims_cannot_be_presented_as_verified(self) -> None:
        root = repository_root()
        payload = json.loads((root / "evidence" / "claims.json").read_text(encoding="utf-8"))
        by_id = {claim["id"]: claim for claim in payload["claims"]}
        self.assertNotEqual(by_id["stable-release"]["status"], "verified")
        self.assertEqual(by_id["public-github-repository"]["status"], "verified")
        self.assertNotEqual(by_id["native-language-quality"]["status"], "verified")
        self.assertNotEqual(by_id["cross-agent-behavior-parity"]["status"], "verified")


if __name__ == "__main__":
    unittest.main()
