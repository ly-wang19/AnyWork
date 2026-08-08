from __future__ import annotations

import json
import unittest

from anywork.catalog import load_catalog, validate_catalog


class CatalogTests(unittest.TestCase):
    def setUp(self) -> None:
        self.catalog = load_catalog()

    def test_catalog_is_valid_and_complete(self) -> None:
        self.assertEqual(validate_catalog(self.catalog), [])
        self.assertEqual(len(self.catalog.skills), 25)
        self.assertEqual(len(self.catalog.packs), 7)

    def test_complete_pack_has_every_skill_once(self) -> None:
        resolved = self.catalog.resolve_pack("complete")
        self.assertEqual(len(resolved), len(set(resolved)))
        self.assertEqual(set(resolved), set(self.catalog.skills))

    def test_experimental_language_quality_is_explicit(self) -> None:
        for skill in self.catalog.skills.values():
            self.assertEqual(skill["maturity"], "experimental")
            self.assertTrue(skill["capabilities"])
            self.assertEqual(set(skill["language_reviews"]), {"en", "zh-CN", "ja"})
            for review in skill["language_reviews"].values():
                self.assertEqual(review["status"], "machine-drafted")
                self.assertEqual(review["revision"], skill["version"])

    def test_all_eval_cases_are_trilingual_and_reference_a_skill(self) -> None:
        path = self.catalog.root / "evals" / "cases.json"
        payload = json.loads(path.read_text(encoding="utf-8"))
        cases = payload["cases"]
        self.assertEqual(len(cases), len(self.catalog.skills) * 3)
        for case in cases:
            self.assertIn(case["skill"], self.catalog.skills)
            self.assertEqual(set(case["prompts"]), {"en", "zh-CN", "ja"})
            self.assertTrue(case["assertions"])


if __name__ == "__main__":
    unittest.main()
