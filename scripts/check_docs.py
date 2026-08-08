#!/usr/bin/env python3
"""Dependency-free structural checks for repository Markdown."""

from __future__ import annotations

import json
import re
import sys
from pathlib import Path
from urllib.parse import unquote


ROOT = Path(__file__).resolve().parents[1]
CAPABILITY_DOCS = (
    Path("CAPABILITIES.md"),
    Path("CAPABILITIES.zh-CN.md"),
    Path("CAPABILITIES.ja.md"),
)
README_CAPABILITY_LINKS = {
    Path("README.md"): "(CAPABILITIES.md)",
    Path("README.zh-CN.md"): "(CAPABILITIES.zh-CN.md)",
    Path("README.ja.md"): "(CAPABILITIES.ja.md)",
}
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


def check_capability_inventory() -> list[str]:
    catalog = json.loads((ROOT / "registry/catalog.json").read_text(encoding="utf-8"))
    skill_ids = [entry["id"] for entry in catalog["skills"]]
    pack_ids = [entry["id"] for entry in catalog["packs"]]
    errors: list[str] = []

    pack_by_id = {entry["id"]: entry for entry in catalog["packs"]}

    def resolve_pack(pack_id: str, trail: tuple[str, ...] = ()) -> set[str]:
        if pack_id in trail:
            return set()
        pack = pack_by_id[pack_id]
        resolved = set(pack["skills"])
        for parent in pack.get("extends", []):
            resolved.update(resolve_pack(parent, trail + (pack_id,)))
        return resolved

    for relative in CAPABILITY_DOCS:
        text = (ROOT / relative).read_text(encoding="utf-8")
        for skill_id in skill_ids:
            marker = f"(skills/{skill_id}/SKILL.md)"
            count = text.count(marker)
            if count != 1:
                errors.append(
                    f"{relative}: expected one linked entry for Skill {skill_id}, found {count}"
                )
        for pack_id in pack_ids:
            marker = f"| `{pack_id}` |"
            row_count = text.count(marker)
            if row_count != 1:
                errors.append(
                    f"{relative}: expected one table row for pack {pack_id}, found {row_count}"
                )
                continue
            row = next(line for line in text.splitlines() if line.startswith(marker))
            cells = [cell.strip() for cell in row.strip().strip("|").split("|")]
            expected_count = len(resolve_pack(pack_id))
            try:
                documented_count = int(cells[2])
            except (IndexError, ValueError):
                errors.append(f"{relative}: invalid Skill count for pack {pack_id}")
                continue
            if documented_count != expected_count:
                errors.append(
                    f"{relative}: pack {pack_id} documents {documented_count} Skills, "
                    f"but resolves to {expected_count}"
                )

    for relative, marker in README_CAPABILITY_LINKS.items():
        text = (ROOT / relative).read_text(encoding="utf-8")
        if marker not in text:
            errors.append(f"{relative}: missing capability-map link {marker}")
    return errors


def main() -> int:
    files = markdown_files()
    errors = [error for path in files for error in check_file(path)]
    errors.extend(check_capability_inventory())
    if errors:
        print("Documentation validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print(
        f"Validated {len(files)} Markdown files: formatting, local links, "
        "and trilingual capability inventory passed."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
