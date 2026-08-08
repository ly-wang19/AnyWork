from __future__ import annotations

import copy
import importlib.util
import json
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location(
    "validate_evals",
    ROOT / "scripts" / "validate_evals.py",
)
assert SPEC is not None and SPEC.loader is not None
validate_evals = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(validate_evals)


class EvalSuiteTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.suite = json.loads((ROOT / "evals" / "cases.json").read_text(encoding="utf-8"))
        cls.catalog = json.loads((ROOT / "registry" / "catalog.json").read_text(encoding="utf-8"))

    def validate(self, suite: dict) -> list[str]:
        return validate_evals.validate_suite(suite, self.catalog)

    def test_complete_suite_has_75_cases_and_225_prompt_variants(self) -> None:
        self.assertEqual(self.validate(self.suite), [])
        self.assertEqual(len(self.suite["cases"]), 75)
        self.assertEqual(len(self.suite["cases"]) * len(self.suite["languages"]), 225)

    def test_each_skill_has_all_three_difficulties(self) -> None:
        expected = {"basic", "ambiguous", "conflicting"}
        by_skill: dict[str, set[str]] = {}
        for case in self.suite["cases"]:
            by_skill.setdefault(case["skill"], set()).add(case["difficulty"])
        self.assertEqual(set(by_skill), {skill["id"] for skill in self.catalog["skills"]})
        self.assertTrue(all(difficulties == expected for difficulties in by_skill.values()))

    def test_manifest_is_reproducible_and_disclaims_runtime_scores(self) -> None:
        actual = json.loads((ROOT / "evals" / "manifest.json").read_text(encoding="utf-8"))
        expected = validate_evals.build_manifest(ROOT, self.suite)
        self.assertEqual(actual, expected)
        self.assertEqual(actual["counts"], {"skills": 25, "cases": 75, "prompt_variants": 225})
        self.assertFalse(actual["claims"]["model_runs_included"])
        self.assertFalse(actual["claims"]["native_human_review_included"])

    def test_suite_digest_is_independent_of_object_key_order(self) -> None:
        reordered = dict(reversed(list(copy.deepcopy(self.suite).items())))
        self.assertEqual(
            validate_evals.suite_digest(self.suite),
            validate_evals.suite_digest(reordered),
        )

    def test_manifest_artifact_hashes_are_independent_of_checkout_newlines(self) -> None:
        expected = validate_evals.build_manifest(ROOT, self.suite)
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            (root / "evals").mkdir()
            (root / "schemas").mkdir()
            for relative in (Path("evals/cases.json"), Path("schemas/eval.schema.json")):
                source = ROOT / relative
                destination = root / relative
                destination.write_bytes(source.read_bytes().replace(b"\n", b"\r\n"))
            actual = validate_evals.build_manifest(root, self.suite)
        self.assertEqual(actual["artifacts"], expected["artifacts"])

    def test_duplicate_id_is_rejected(self) -> None:
        candidate = copy.deepcopy(self.suite)
        candidate["cases"][1]["id"] = candidate["cases"][0]["id"]
        problems = self.validate(candidate)
        self.assertTrue(any("duplicate ids" in problem for problem in problems), problems)

    def test_missing_locale_is_rejected(self) -> None:
        candidate = copy.deepcopy(self.suite)
        del candidate["cases"][0]["prompts"]["ja"]
        problems = self.validate(candidate)
        self.assertTrue(any("prompts" in problem and "expected exactly" in problem for problem in problems), problems)

    def test_neighbor_must_be_registered_and_in_not_for(self) -> None:
        candidate = copy.deepcopy(self.suite)
        candidate["cases"][0]["routing"]["neighbor"] = "missing-skill"
        problems = self.validate(candidate)
        self.assertTrue(any("routing.neighbor" in problem for problem in problems), problems)
        self.assertTrue(any("must include neighbor" in problem for problem in problems), problems)

    def test_hard_negative_trigger_must_appear_in_prompt(self) -> None:
        candidate = copy.deepcopy(self.suite)
        ambiguous = next(case for case in candidate["cases"] if case["difficulty"] == "ambiguous")
        ambiguous["routing"]["negative_triggers"]["en"] = ["phrase absent from prompt"]
        problems = self.validate(candidate)
        self.assertTrue(any("not found in prompt" in problem for problem in problems), problems)

    def test_answer_bearing_fields_are_rejected(self) -> None:
        candidate = copy.deepcopy(self.suite)
        candidate["cases"][0]["expected_output"] = "future model answer"
        problems = self.validate(candidate)
        self.assertTrue(any("answer-bearing fields are forbidden" in problem for problem in problems), problems)

    def test_unknown_fields_are_rejected(self) -> None:
        candidate = copy.deepcopy(self.suite)
        candidate["cases"][0]["routing"]["notes"] = "not in the schema"
        problems = self.validate(candidate)
        self.assertTrue(any("routing: unexpected keys" in problem for problem in problems), problems)

    def test_malformed_nested_values_report_errors_instead_of_crashing(self) -> None:
        candidate = copy.deepcopy(self.suite)
        candidate["cases"][0]["routing"]["not_for"] = [{"bad": "shape"}]
        candidate["cases"][1]["routing"]["negative_triggers"]["en"] = [{"bad": "shape"}]
        candidate["cases"][2]["assertions"]["must"] = [{"bad": "shape"}]
        problems = self.validate(candidate)
        self.assertTrue(any("routing.not_for" in problem for problem in problems), problems)
        self.assertTrue(any("negative_triggers.en" in problem for problem in problems), problems)
        self.assertTrue(any("assertions.must" in problem for problem in problems), problems)


if __name__ == "__main__":
    unittest.main()
