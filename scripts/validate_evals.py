#!/usr/bin/env python3
"""Validate AnyWork's static multilingual evaluation suite and manifest."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from collections import Counter, defaultdict
from pathlib import Path
from typing import Any


LANGUAGES = ("en", "zh-CN", "ja")
DIFFICULTIES = ("basic", "ambiguous", "conflicting")
RUBRIC = (
    "task-correctness",
    "traceability",
    "direct-usability",
    "native-register",
    "safety",
)
TOP_LEVEL_KEYS = {"schema_version", "suite_id", "languages", "difficulties", "rubrics", "cases"}
CASE_KEYS = {"id", "skill", "difficulty", "prompts", "routing", "assertions", "rubric"}
ROUTING_KEYS = {"expected_skill", "neighbor", "not_for", "hard_negative", "negative_triggers"}
ASSERTION_KEYS = {"must", "must_not"}
FORBIDDEN_ANSWER_KEYS = {"answer", "expected_answer", "expected_output", "gold_answer", "model_output"}
SKILL_ID = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")


def repository_root() -> Path:
    return Path(__file__).resolve().parents[1]


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def canonical_json(payload: Any) -> bytes:
    return json.dumps(
        payload,
        ensure_ascii=False,
        sort_keys=True,
        separators=(",", ":"),
    ).encode("utf-8")


def sha256_bytes(payload: bytes) -> str:
    return hashlib.sha256(payload).hexdigest()


def normalized_text_bytes(path: Path) -> bytes:
    """Return canonical text bytes independent of Git checkout newline mode."""

    return path.read_bytes().replace(b"\r\n", b"\n")


def unique_strings(value: Any) -> bool:
    return (
        isinstance(value, list)
        and all(isinstance(item, str) for item in value)
        and len(value) == len(set(value))
    )


def suite_digest(payload: dict[str, Any]) -> str:
    return sha256_bytes(canonical_json(payload))


def _unknown_keys(value: Any, path: str, problems: list[str]) -> None:
    if isinstance(value, dict):
        forbidden = sorted(FORBIDDEN_ANSWER_KEYS.intersection(value))
        if forbidden:
            problems.append(f"{path}: answer-bearing fields are forbidden: {', '.join(forbidden)}")
        for key, child in value.items():
            _unknown_keys(child, f"{path}.{key}", problems)
    elif isinstance(value, list):
        for index, child in enumerate(value):
            _unknown_keys(child, f"{path}[{index}]", problems)


def validate_suite(payload: Any, catalog: Any) -> list[str]:
    """Return deterministic structural and semantic validation errors."""

    problems: list[str] = []
    if not isinstance(payload, dict):
        return ["suite: expected an object"]
    unexpected = sorted(set(payload) - TOP_LEVEL_KEYS)
    missing = sorted(TOP_LEVEL_KEYS - set(payload))
    if unexpected:
        problems.append(f"suite: unexpected keys: {', '.join(unexpected)}")
    if missing:
        problems.append(f"suite: missing keys: {', '.join(missing)}")
    if payload.get("schema_version") != 2:
        problems.append("suite.schema_version: expected 2")
    if payload.get("suite_id") != "anywork-routing-v2":
        problems.append("suite.suite_id: expected anywork-routing-v2")
    if payload.get("languages") != list(LANGUAGES):
        problems.append("suite.languages: expected en, zh-CN, ja in canonical order")
    if payload.get("difficulties") != list(DIFFICULTIES):
        problems.append("suite.difficulties: expected basic, ambiguous, conflicting in canonical order")
    if payload.get("rubrics") != {"standard-v1": list(RUBRIC)}:
        problems.append("suite.rubrics: standard-v1 does not match the canonical dimensions")

    catalog_skills = {
        entry.get("id")
        for entry in catalog.get("skills", [])
        if isinstance(entry, dict) and isinstance(entry.get("id"), str)
    }
    cases = payload.get("cases")
    if not isinstance(cases, list):
        problems.append("suite.cases: expected an array")
        return problems
    if len(cases) != len(catalog_skills) * len(DIFFICULTIES):
        problems.append(
            f"suite.cases: expected {len(catalog_skills) * len(DIFFICULTIES)} cases, got {len(cases)}"
        )

    ids: Counter[str] = Counter()
    coverage: dict[str, Counter[str]] = defaultdict(Counter)
    for index, case in enumerate(cases):
        case_path = f"suite.cases[{index}]"
        if not isinstance(case, dict):
            problems.append(f"{case_path}: expected an object")
            continue
        unexpected_case = sorted(set(case) - CASE_KEYS)
        missing_case = sorted(CASE_KEYS - set(case))
        if unexpected_case:
            problems.append(f"{case_path}: unexpected keys: {', '.join(unexpected_case)}")
        if missing_case:
            problems.append(f"{case_path}: missing keys: {', '.join(missing_case)}")

        case_id = case.get("id")
        skill = case.get("skill")
        difficulty = case.get("difficulty")
        if not isinstance(case_id, str):
            problems.append(f"{case_path}.id: expected a string")
        else:
            ids[case_id] += 1
            expected_id = f"{skill}-{difficulty}-001"
            if case_id != expected_id:
                problems.append(f"{case_path}.id: expected {expected_id}")
        if skill not in catalog_skills:
            problems.append(f"{case_path}.skill: unknown skill {skill!r}")
        if difficulty not in DIFFICULTIES:
            problems.append(f"{case_path}.difficulty: invalid difficulty {difficulty!r}")
        elif isinstance(skill, str):
            coverage[skill][difficulty] += 1
        if case.get("rubric") != "standard-v1":
            problems.append(f"{case_path}.rubric: expected standard-v1")

        prompts = case.get("prompts")
        if not isinstance(prompts, dict):
            problems.append(f"{case_path}.prompts: expected an object")
            prompts = {}
        elif set(prompts) != set(LANGUAGES):
            problems.append(f"{case_path}.prompts: expected exactly en, zh-CN, and ja")
        for language in LANGUAGES:
            prompt = prompts.get(language)
            if not isinstance(prompt, str) or not prompt.strip():
                problems.append(f"{case_path}.prompts.{language}: expected non-empty text")

        routing = case.get("routing")
        if not isinstance(routing, dict):
            problems.append(f"{case_path}.routing: expected an object")
            routing = {}
        else:
            extra = sorted(set(routing) - ROUTING_KEYS)
            absent = sorted(ROUTING_KEYS - set(routing))
            if extra:
                problems.append(f"{case_path}.routing: unexpected keys: {', '.join(extra)}")
            if absent:
                problems.append(f"{case_path}.routing: missing keys: {', '.join(absent)}")
        if routing.get("expected_skill") != skill:
            problems.append(f"{case_path}.routing.expected_skill: must equal case.skill")
        neighbor = routing.get("neighbor")
        if neighbor not in catalog_skills or neighbor == skill:
            problems.append(f"{case_path}.routing.neighbor: expected a different registered skill")
        not_for = routing.get("not_for")
        if not unique_strings(not_for) or not not_for or any(
            item not in catalog_skills or item == skill for item in not_for
        ):
            problems.append(f"{case_path}.routing.not_for: expected unique, different registered skills")
        elif neighbor not in not_for:
            problems.append(f"{case_path}.routing.not_for: must include neighbor")

        hard_negative = routing.get("hard_negative")
        expected_hard_negative = difficulty in {"ambiguous", "conflicting"}
        if hard_negative is not expected_hard_negative:
            problems.append(
                f"{case_path}.routing.hard_negative: expected {str(expected_hard_negative).lower()}"
            )
        negative_triggers = routing.get("negative_triggers")
        if not isinstance(negative_triggers, dict):
            problems.append(f"{case_path}.routing.negative_triggers: expected an object")
            negative_triggers = {}
        elif set(negative_triggers) != set(LANGUAGES):
            problems.append(f"{case_path}.routing.negative_triggers: expected exactly three languages")
        for language in LANGUAGES:
            triggers = negative_triggers.get(language)
            if not unique_strings(triggers) or any(not item for item in triggers):
                problems.append(f"{case_path}.routing.negative_triggers.{language}: expected strings")
                continue
            if expected_hard_negative and not triggers:
                problems.append(f"{case_path}.routing.negative_triggers.{language}: hard negative needs a trigger")
            if not expected_hard_negative and triggers:
                problems.append(f"{case_path}.routing.negative_triggers.{language}: basic case must be empty")
            prompt = prompts.get(language, "")
            if isinstance(prompt, str):
                for trigger in triggers:
                    if trigger.casefold() not in prompt.casefold():
                        problems.append(
                            f"{case_path}.routing.negative_triggers.{language}: {trigger!r} not found in prompt"
                        )

        assertions = case.get("assertions")
        if not isinstance(assertions, dict):
            problems.append(f"{case_path}.assertions: expected an object")
            assertions = {}
        else:
            extra = sorted(set(assertions) - ASSERTION_KEYS)
            absent = sorted(ASSERTION_KEYS - set(assertions))
            if extra:
                problems.append(f"{case_path}.assertions: unexpected keys: {', '.join(extra)}")
            if absent:
                problems.append(f"{case_path}.assertions: missing keys: {', '.join(absent)}")
        must = assertions.get("must")
        must_not = assertions.get("must_not")
        if (
            not unique_strings(must)
            or len(must) < 3
            or any(not SKILL_ID.fullmatch(item) for item in must)
        ):
            problems.append(f"{case_path}.assertions.must: expected at least three unique contract ids")
        if (
            not unique_strings(must_not)
            or any(not SKILL_ID.fullmatch(item) for item in must_not)
        ):
            problems.append(f"{case_path}.assertions.must_not: expected unique contract ids")
        elif expected_hard_negative and not must_not:
            problems.append(f"{case_path}.assertions.must_not: hard negative needs a boundary contract")

    duplicate_ids = sorted(case_id for case_id, count in ids.items() if count != 1)
    if duplicate_ids:
        problems.append(f"suite.cases: duplicate ids: {', '.join(duplicate_ids)}")
    for skill in sorted(catalog_skills):
        actual = coverage.get(skill, Counter())
        for difficulty in DIFFICULTIES:
            if actual[difficulty] != 1:
                problems.append(
                    f"suite.coverage.{skill}.{difficulty}: expected 1 case, got {actual[difficulty]}"
                )

    _unknown_keys(payload, "suite", problems)
    return problems


def build_manifest(root: Path, suite: dict[str, Any]) -> dict[str, Any]:
    cases_path = root / "evals" / "cases.json"
    schema_path = root / "schemas" / "eval.schema.json"
    skills = {case["skill"] for case in suite["cases"]}
    return {
        "schema_version": 1,
        "suite_id": suite["suite_id"],
        "digest_algorithm": "sha256",
        "artifact_text_normalization": "crlf-to-lf",
        "suite_digest": suite_digest(suite),
        "artifacts": {
            "cases": {
                "path": "evals/cases.json",
                "sha256": sha256_bytes(normalized_text_bytes(cases_path)),
            },
            "schema": {
                "path": "schemas/eval.schema.json",
                "sha256": sha256_bytes(normalized_text_bytes(schema_path)),
            },
        },
        "counts": {
            "skills": len(skills),
            "cases": len(suite["cases"]),
            "prompt_variants": len(suite["cases"]) * len(suite["languages"]),
        },
        "languages": suite["languages"],
        "difficulties": suite["difficulties"],
        "rubric_ids": sorted(suite["rubrics"]),
        "claims": {
            "model_runs_included": False,
            "native_human_review_included": False,
        },
    }


def validate_manifest(root: Path, suite: dict[str, Any]) -> list[str]:
    path = root / "evals" / "manifest.json"
    if not path.is_file():
        return ["evals/manifest.json: missing; run scripts/validate_evals.py --write-manifest"]
    try:
        actual = load_json(path)
    except (OSError, json.JSONDecodeError) as error:
        return [f"evals/manifest.json: invalid JSON: {error}"]
    expected = build_manifest(root, suite)
    if actual != expected:
        return ["evals/manifest.json: stale or non-reproducible; regenerate with --write-manifest"]
    return []


def parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", type=Path, default=repository_root())
    parser.add_argument("--write-manifest", action="store_true")
    return parser.parse_args(argv)


def main(argv: list[str] | None = None) -> int:
    args = parse_args(argv)
    root = args.root.resolve()
    try:
        suite = load_json(root / "evals" / "cases.json")
        catalog = load_json(root / "registry" / "catalog.json")
    except (OSError, json.JSONDecodeError) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        return 2
    problems = validate_suite(suite, catalog)
    if problems:
        for problem in problems:
            print(f"ERROR: {problem}", file=sys.stderr)
        return 1
    if args.write_manifest:
        manifest = build_manifest(root, suite)
        target = root / "evals" / "manifest.json"
        target.write_text(json.dumps(manifest, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    problems = validate_manifest(root, suite)
    if problems:
        for problem in problems:
            print(f"ERROR: {problem}", file=sys.stderr)
        return 1
    print(
        f"Validated {len(suite['cases'])} cases, "
        f"{len(suite['cases']) * len(suite['languages'])} prompt variants, "
        f"digest {suite_digest(suite)}"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
