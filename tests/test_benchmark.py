from __future__ import annotations

import importlib.util
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
SPEC = importlib.util.spec_from_file_location("benchmark", ROOT / "scripts" / "benchmark.py")
assert SPEC is not None and SPEC.loader is not None
benchmark = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(benchmark)


class BenchmarkTests(unittest.TestCase):
    def test_command_requires_one_prompt_placeholder_and_does_not_use_a_shell(self) -> None:
        command = benchmark.command_for_prompt('runner --input "{prompt}"', "hello; rm -rf example")
        self.assertEqual(command, ["runner", "--input", "hello; rm -rf example"])
        with self.assertRaises(ValueError):
            benchmark.command_for_prompt("runner", "prompt")

    def test_compare_requires_identical_prompts_and_emits_blank_scorecard(self) -> None:
        result = {
            "case_id": "example-basic-001",
            "skill": "example",
            "difficulty": "basic",
            "prompt": "same prompt",
            "stdout": "answer",
            "exit_code": 0,
            "timed_out": False,
            "duration_ms": 1,
        }
        shared = {"suite_id": "suite", "suite_digest": "digest", "language": "en"}
        baseline = {**shared, "mode": "baseline", "results": [result]}
        anywork = {**shared, "mode": "anywork", "results": [{**result, "stdout": "better answer"}]}
        report = benchmark.render_review(baseline, anywork)
        self.assertIn("byte-identical prompts", report)
        self.assertIn("| task-correctness |  |  |  |", report)
        changed = {**anywork, "results": [{**result, "prompt": "changed"}]}
        with self.assertRaises(ValueError):
            benchmark.render_review(baseline, changed)


if __name__ == "__main__":
    unittest.main()
