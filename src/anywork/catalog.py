from __future__ import annotations

import json
import re
from dataclasses import dataclass
from pathlib import Path
from typing import Any


SKILL_NAME = re.compile(r"^[a-z0-9]+(?:-[a-z0-9]+)*$")
REQUIRED_LANGUAGES = {"en", "zh-CN", "ja"}
REQUIRED_FRONTMATTER = {"name", "description"}
DANGEROUS_PATTERNS = {
    re.compile(r"curl\b[^\n|]*\|\s*(?:ba)?sh\b", re.IGNORECASE): "remote pipe-to-shell command",
    re.compile(r"wget\b[^\n|]*\|\s*(?:ba)?sh\b", re.IGNORECASE): "remote pipe-to-shell command",
    re.compile(r"\brm\s+-[a-z]*r[a-z]*f\b", re.IGNORECASE): "recursive forced deletion",
    re.compile(r"(?:~|\$HOME)/(?:\.ssh|\.aws|\.config/gcloud)\b"): "credential-directory access",
    re.compile(r"BEGIN (?:RSA|OPENSSH|EC) PRIVATE KEY"): "embedded private key",
}


@dataclass(frozen=True)
class Catalog:
    root: Path
    raw: dict[str, Any]

    @property
    def skills(self) -> dict[str, dict[str, Any]]:
        return {entry["id"]: entry for entry in self.raw.get("skills", [])}

    @property
    def packs(self) -> dict[str, dict[str, Any]]:
        return {entry["id"]: entry for entry in self.raw.get("packs", [])}

    @property
    def agents(self) -> tuple[str, ...]:
        return tuple(self.raw.get("agents", []))

    def skill_path(self, skill_id: str) -> Path:
        entry = self.skills[skill_id]
        return (self.root / entry["path"]).resolve()

    def resolve_pack(self, pack_id: str) -> list[str]:
        seen_packs: set[str] = set()
        visiting: set[str] = set()
        seen_skills: set[str] = set()
        ordered: list[str] = []

        def add_pack(current: str) -> None:
            if current in seen_packs:
                return
            if current in visiting:
                raise ValueError(f"cyclic pack inheritance at {current}")
            if current not in self.packs:
                raise KeyError(current)
            visiting.add(current)
            pack = self.packs[current]
            for parent in pack.get("extends", []):
                add_pack(parent)
            for skill in pack.get("skills", []):
                if skill not in seen_skills:
                    seen_skills.add(skill)
                    ordered.append(skill)
            visiting.remove(current)
            seen_packs.add(current)

        add_pack(pack_id)
        return ordered


def repository_root() -> Path:
    return Path(__file__).resolve().parents[2]


def load_catalog(root: Path | None = None) -> Catalog:
    project_root = (root or repository_root()).resolve()
    catalog_path = project_root / "registry" / "catalog.json"
    raw = json.loads(catalog_path.read_text(encoding="utf-8"))
    return Catalog(root=project_root, raw=raw)


def parse_frontmatter(path: Path) -> tuple[dict[str, str], list[str]]:
    problems: list[str] = []
    lines = path.read_text(encoding="utf-8").splitlines()
    if not lines or lines[0].strip() != "---":
        return {}, [f"{path}: missing YAML frontmatter"]
    try:
        end = lines.index("---", 1)
    except ValueError:
        return {}, [f"{path}: unclosed YAML frontmatter"]
    values: dict[str, str] = {}
    for line in lines[1:end]:
        if not line.strip() or line.lstrip().startswith("#"):
            continue
        if ":" not in line:
            problems.append(f"{path}: invalid frontmatter line: {line}")
            continue
        key, value = line.split(":", 1)
        values[key.strip()] = value.strip().strip('"').strip("'")
    if set(values) != REQUIRED_FRONTMATTER:
        problems.append(
            f"{path}: frontmatter keys must be exactly {sorted(REQUIRED_FRONTMATTER)}, got {sorted(values)}"
        )
    return values, problems


def validate_catalog(catalog: Catalog) -> list[str]:
    problems: list[str] = []
    if catalog.raw.get("schema_version") != 1:
        problems.append("registry/catalog.json: schema_version must be 1")
    if set(catalog.raw.get("languages", [])) != REQUIRED_LANGUAGES:
        problems.append("registry/catalog.json: languages must be en, zh-CN, and ja")
    if not catalog.agents:
        problems.append("registry/catalog.json: at least one agent is required")

    skills = catalog.skills
    packs = catalog.packs
    if len(skills) != len(catalog.raw.get("skills", [])):
        problems.append("registry/catalog.json: duplicate skill id")
    if len(packs) != len(catalog.raw.get("packs", [])):
        problems.append("registry/catalog.json: duplicate pack id")

    for skill_id, entry in skills.items():
        if not SKILL_NAME.fullmatch(skill_id):
            problems.append(f"skill {skill_id}: id must use lowercase kebab-case")
        if set(entry.get("languages", [])) != REQUIRED_LANGUAGES:
            problems.append(f"skill {skill_id}: must support en, zh-CN, and ja")
        if not entry.get("license") or not entry.get("source"):
            problems.append(f"skill {skill_id}: license and source are required")
        if entry.get("risk") not in {"low", "medium", "high"}:
            problems.append(f"skill {skill_id}: risk must be low, medium, or high")
        path = catalog.skill_path(skill_id)
        if catalog.root not in path.parents:
            problems.append(f"skill {skill_id}: path escapes repository root")
            continue
        skill_file = path / "SKILL.md"
        if not skill_file.is_file():
            problems.append(f"skill {skill_id}: missing {skill_file}")
            continue
        frontmatter, fm_problems = parse_frontmatter(skill_file)
        problems.extend(fm_problems)
        if frontmatter.get("name") != skill_id:
            problems.append(f"skill {skill_id}: SKILL.md name does not match directory/id")
        description = frontmatter.get("description", "")
        if len(description) > 600:
            problems.append(f"skill {skill_id}: description exceeds AnyWork's 600-character budget")
        if not re.search(r"[\u4e00-\u9fff]", description):
            problems.append(f"skill {skill_id}: description lacks a Chinese trigger phrase")
        if not re.search(r"[\u3040-\u30ff]", description):
            problems.append(f"skill {skill_id}: description lacks a Japanese trigger phrase")
        if len(skill_file.read_text(encoding="utf-8").splitlines()) > 500:
            problems.append(f"skill {skill_id}: SKILL.md exceeds 500 lines")
        if not (path / "agents" / "openai.yaml").is_file():
            problems.append(f"skill {skill_id}: missing Codex agents/openai.yaml metadata")
        for child in path.rglob("*"):
            if child.is_symlink():
                problems.append(f"skill {skill_id}: symlinks are not allowed in official packages: {child}")
            if child.is_file():
                content = child.read_text(encoding="utf-8", errors="replace")
                for pattern, label in DANGEROUS_PATTERNS.items():
                    if pattern.search(content):
                        problems.append(f"skill {skill_id}: detected {label} in {child.relative_to(catalog.root)}")

    for pack_id, pack in packs.items():
        if not SKILL_NAME.fullmatch(pack_id):
            problems.append(f"pack {pack_id}: id must use lowercase kebab-case")
        if set(pack.get("display_name", {})) != REQUIRED_LANGUAGES:
            problems.append(f"pack {pack_id}: display_name must have all three languages")
        if set(pack.get("description", {})) != REQUIRED_LANGUAGES:
            problems.append(f"pack {pack_id}: description must have all three languages")
        for skill in pack.get("skills", []):
            if skill not in skills:
                problems.append(f"pack {pack_id}: unknown skill {skill}")
        for parent in pack.get("extends", []):
            if parent not in packs:
                problems.append(f"pack {pack_id}: unknown parent pack {parent}")
        try:
            catalog.resolve_pack(pack_id)
        except (RecursionError, ValueError):
            problems.append(f"pack {pack_id}: cyclic inheritance")
        except KeyError as error:
            problems.append(f"pack {pack_id}: missing dependency {error.args[0]}")

    skill_root = catalog.root / "skills"
    listed_paths = {catalog.skill_path(skill_id) for skill_id in skills}
    actual_paths = {path.resolve() for path in skill_root.iterdir() if path.is_dir()}
    for extra in sorted(actual_paths - listed_paths):
        problems.append(f"unregistered skill directory: {extra.relative_to(catalog.root)}")

    eval_path = catalog.root / "evals" / "cases.json"
    if not eval_path.is_file():
        problems.append("evals/cases.json: missing multilingual evaluation suite")
    else:
        try:
            eval_payload = json.loads(eval_path.read_text(encoding="utf-8"))
            eval_cases = eval_payload.get("cases", [])
            covered: set[str] = set()
            for case in eval_cases:
                skill = case.get("skill")
                if skill not in skills:
                    problems.append(f"eval {case.get('id', '<unknown>')}: unknown skill {skill}")
                else:
                    covered.add(skill)
                if set(case.get("prompts", {})) != REQUIRED_LANGUAGES:
                    problems.append(f"eval {case.get('id', '<unknown>')}: prompts must cover all three languages")
                if not case.get("assertions") or not case.get("rubric"):
                    problems.append(f"eval {case.get('id', '<unknown>')}: assertions and rubric are required")
            for missing in sorted(set(skills) - covered):
                problems.append(f"skill {missing}: no multilingual evaluation case")
        except (json.JSONDecodeError, OSError, AttributeError) as error:
            problems.append(f"evals/cases.json: invalid evaluation suite: {error}")
    return problems
