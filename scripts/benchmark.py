#!/usr/bin/env python3
"""Run and compare reproducible before/after AnyWork model evaluations."""

from __future__ import annotations

import argparse
import datetime as dt
import hashlib
import json
import shlex
import subprocess
import sys
import time
from pathlib import Path
from typing import Any


ROOT = Path(__file__).resolve().parents[1]
LANGUAGES = ("en", "zh-CN", "ja")
MODES = ("baseline", "anywork")
RUBRIC = (
    "task-correctness",
    "traceability",
    "direct-usability",
    "native-register",
    "safety",
)


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def digest(value: Any) -> str:
    payload = json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":"))
    return hashlib.sha256(payload.encode("utf-8")).hexdigest()


def command_for_prompt(template: str, prompt: str) -> list[str]:
    if template.count("{prompt}") != 1:
        raise ValueError("--command must contain exactly one {prompt} placeholder")
    arguments = shlex.split(template)
    return [argument.replace("{prompt}", prompt) for argument in arguments]


def select_cases(
    suite: dict[str, Any],
    language: str,
    skills: set[str],
    difficulties: set[str],
    case_ids: set[str],
    limit: int | None,
) -> list[dict[str, Any]]:
    selected = []
    for case in suite["cases"]:
        if skills and case["skill"] not in skills:
            continue
        if difficulties and case["difficulty"] not in difficulties:
            continue
        if case_ids and case["id"] not in case_ids:
            continue
        selected.append(
            {
                "case_id": case["id"],
                "skill": case["skill"],
                "difficulty": case["difficulty"],
                "prompt": case["prompts"][language],
                "rubric": case["rubric"],
            }
        )
    return selected[:limit] if limit is not None else selected


def run_case(command_template: str, case: dict[str, Any], timeout: float) -> dict[str, Any]:
    command = command_for_prompt(command_template, case["prompt"])
    started = time.monotonic()
    try:
        completed = subprocess.run(
            command,
            capture_output=True,
            check=False,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=timeout,
        )
        exit_code: int | None = completed.returncode
        stdout = completed.stdout
        stderr = completed.stderr
        timed_out = False
    except subprocess.TimeoutExpired as error:
        exit_code = None
        stdout = error.stdout or ""
        stderr = error.stderr or ""
        if isinstance(stdout, bytes):
            stdout = stdout.decode("utf-8", errors="replace")
        if isinstance(stderr, bytes):
            stderr = stderr.decode("utf-8", errors="replace")
        timed_out = True
    return {
        **case,
        "exit_code": exit_code,
        "timed_out": timed_out,
        "duration_ms": round((time.monotonic() - started) * 1000),
        "stdout": stdout,
        "stderr": stderr,
    }


def write_new(path: Path, content: str, force: bool) -> None:
    if path.exists() and not force:
        raise FileExistsError(f"refusing to overwrite {path}; pass --force after review")
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")


def run_benchmark(args: argparse.Namespace) -> int:
    suite = load_json(ROOT / "evals/cases.json")
    selected = select_cases(
        suite,
        args.language,
        set(args.skill),
        set(args.difficulty),
        set(args.case),
        args.limit,
    )
    if not selected:
        raise ValueError("no evaluation cases matched the selected filters")
    command_for_prompt(args.command, selected[0]["prompt"])
    results = []
    for index, case in enumerate(selected, start=1):
        print(f"[{index}/{len(selected)}] {case['case_id']}", file=sys.stderr)
        results.append(run_case(args.command, case, args.timeout))
    run = {
        "schema_version": 1,
        "created_at": dt.datetime.now(dt.timezone.utc).isoformat(),
        "mode": args.mode,
        "suite_id": suite["suite_id"],
        "suite_digest": digest(suite),
        "language": args.language,
        "command_template": args.command,
        "claims": {
            "automatic_quality_score": False,
            "prompt_modified_for_anywork": False,
        },
        "results": results,
    }
    write_new(args.output, json.dumps(run, ensure_ascii=False, indent=2) + "\n", args.force)
    print(f"Wrote {len(results)} exact-prompt results to {args.output}")
    return 0


def fenced(value: str) -> str:
    fence = "```"
    while fence in value:
        fence += "`"
    return f"{fence}\n{value.rstrip()}\n{fence}"


def paired_results(baseline: dict[str, Any], anywork: dict[str, Any]) -> list[tuple[dict[str, Any], dict[str, Any]]]:
    if baseline.get("mode") != "baseline" or anywork.get("mode") != "anywork":
        raise ValueError("compare expects a baseline run followed by an anywork run")
    for field in ("suite_id", "suite_digest", "language"):
        if baseline.get(field) != anywork.get(field):
            raise ValueError(f"runs do not match on {field}")
    before = {result["case_id"]: result for result in baseline["results"]}
    after = {result["case_id"]: result for result in anywork["results"]}
    if set(before) != set(after):
        raise ValueError("runs do not contain the same case ids")
    pairs = []
    for result in baseline["results"]:
        candidate = after[result["case_id"]]
        if result["prompt"] != candidate["prompt"]:
            raise ValueError(f"prompt changed for {result['case_id']}")
        pairs.append((result, candidate))
    return pairs


def render_review(baseline: dict[str, Any], anywork: dict[str, Any]) -> str:
    pairs = paired_results(baseline, anywork)
    lines = [
        "# AnyWork before/after review",
        "",
        f"- Suite: `{baseline['suite_id']}`",
        f"- Language: `{baseline['language']}`",
        f"- Paired cases: {len(pairs)}",
        "- Integrity: both runs used byte-identical prompts; AnyWork received no routing hint.",
        "- Scoring: intentionally blank. A qualified reviewer should score each 0–5 dimension before totals are claimed.",
        "",
    ]
    for before, after in pairs:
        lines.extend(
            [
                f"## {before['case_id']}",
                "",
                f"Skill under test: `{before['skill']}` · difficulty: `{before['difficulty']}`",
                "",
                "### Exact prompt",
                "",
                fenced(before["prompt"]),
                "",
                "### Baseline output",
                "",
                fenced(before["stdout"]),
                "",
                f"Exit: `{before['exit_code']}` · timeout: `{before['timed_out']}` · duration: `{before['duration_ms']} ms`",
                "",
                "### AnyWork output",
                "",
                fenced(after["stdout"]),
                "",
                f"Exit: `{after['exit_code']}` · timeout: `{after['timed_out']}` · duration: `{after['duration_ms']} ms`",
                "",
                "### Reviewer scorecard",
                "",
                "| Dimension | Baseline (0–5) | AnyWork (0–5) | Evidence / notes |",
                "|---|---:|---:|---|",
            ]
        )
        lines.extend(f"| {dimension} |  |  |  |" for dimension in RUBRIC)
        lines.extend(["", "Decision: □ AnyWork wins □ tie □ baseline wins □ invalid", ""])
    return "\n".join(lines)


def compare_benchmarks(args: argparse.Namespace) -> int:
    baseline = load_json(args.baseline)
    anywork = load_json(args.anywork)
    report = render_review(baseline, anywork)
    write_new(args.output, report + "\n", args.force)
    print(f"Wrote paired review sheet to {args.output}")
    return 0


def parser() -> argparse.ArgumentParser:
    top = argparse.ArgumentParser(description=__doc__)
    commands = top.add_subparsers(dest="subcommand", required=True)
    run = commands.add_parser("run", help="execute an exact-prompt baseline or AnyWork run")
    run.add_argument("--mode", choices=MODES, required=True)
    run.add_argument("--command", required=True, help="runner command containing exactly one {prompt} placeholder")
    run.add_argument("--output", type=Path, required=True)
    run.add_argument("--language", choices=LANGUAGES, default="en")
    run.add_argument("--skill", action="append", default=[])
    run.add_argument("--difficulty", choices=("basic", "ambiguous", "conflicting"), action="append", default=[])
    run.add_argument("--case", action="append", default=[])
    run.add_argument("--limit", type=int)
    run.add_argument("--timeout", type=float, default=300.0)
    run.add_argument("--force", action="store_true")
    run.set_defaults(handler=run_benchmark)

    compare = commands.add_parser("compare", help="build a human-review sheet from two exact-prompt runs")
    compare.add_argument("baseline", type=Path)
    compare.add_argument("anywork", type=Path)
    compare.add_argument("--output", type=Path, required=True)
    compare.add_argument("--force", action="store_true")
    compare.set_defaults(handler=compare_benchmarks)
    return top


def main(argv: list[str] | None = None) -> int:
    args = parser().parse_args(argv)
    try:
        return args.handler(args)
    except (FileExistsError, OSError, ValueError, json.JSONDecodeError) as error:
        print(f"ERROR: {error}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
