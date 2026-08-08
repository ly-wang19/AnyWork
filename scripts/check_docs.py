#!/usr/bin/env python3
"""Dependency-free structural checks for repository Markdown."""

from __future__ import annotations

import re
import sys
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parents[1]
LINK_RE = re.compile(r"(?<!!)\[[^\]]*\]\(([^)]+)\)")
HEADING_RE = re.compile(r"^(#{1,6})\s+\S")
FENCE_RE = re.compile(r"^\s*(```+|~~~+)")


def markdown_files() -> list[Path]:
    return sorted(path for path in ROOT.rglob("*.md") if ".git" not in path.parts)


def check_file(path: Path) -> list[str]:
    relative = path.relative_to(ROOT)
    raw = path.read_bytes()
    errors: list[str] = []
    if b"\r" in raw:
        errors.append(f"{relative}: contains CR line endings")

    text = raw.decode("utf-8")
    lines = text.splitlines()
    in_fence = False
    fence_marker = ""
    previous_heading = 0
    h1_count = 0

    for number, line in enumerate(lines, start=1):
        if line != line.rstrip():
            errors.append(f"{relative}:{number}: trailing whitespace")

        fence = FENCE_RE.match(line)
        if fence:
            marker = fence.group(1)
            if not in_fence:
                in_fence = True
                fence_marker = marker[0]
            elif marker[0] == fence_marker:
                in_fence = False
                fence_marker = ""
            continue

        if in_fence:
            continue

        heading = HEADING_RE.match(line)
        if heading:
            level = len(heading.group(1))
            if level == 1:
                h1_count += 1
            if previous_heading and level > previous_heading + 1:
                errors.append(
                    f"{relative}:{number}: heading jumps from H{previous_heading} to H{level}"
                )
            previous_heading = level

        for match in LINK_RE.finditer(line):
            destination = match.group(1).strip()
            if destination.startswith("<") and destination.endswith(">"):
                destination = destination[1:-1]
            destination = destination.split(maxsplit=1)[0]
            if destination.startswith(("#", "http://", "https://", "mailto:")):
                continue
            target_text = unquote(destination.split("#", 1)[0])
            if not target_text:
                continue
            target = (path.parent / target_text).resolve()
            try:
                target.relative_to(ROOT)
            except ValueError:
                errors.append(f"{relative}:{number}: link escapes repository: {destination}")
                continue
            if not target.exists():
                errors.append(f"{relative}:{number}: missing link target: {destination}")

    if in_fence:
        errors.append(f"{relative}: unclosed code fence")
    if h1_count != 1:
        errors.append(f"{relative}: expected exactly one H1, found {h1_count}")
    return errors


def main() -> int:
    files = markdown_files()
    errors = [error for path in files for error in check_file(path)]
    if errors:
        print("Documentation validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print(f"Validated {len(files)} Markdown files: formatting and local links passed.")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
