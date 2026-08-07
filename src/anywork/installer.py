from __future__ import annotations

import hashlib
import json
import os
import shutil
import tempfile
from datetime import datetime, timezone
from pathlib import Path
from typing import Callable

from .catalog import Catalog


SUPPORTED_AGENTS = ("codex", "claude-code")


class InstallError(RuntimeError):
    pass


def normalize_agent(agent: str) -> str:
    aliases = {"claude": "claude-code", "claude_code": "claude-code"}
    return aliases.get(agent, agent)


def target_root(agent: str, scope: str, project_dir: Path | None = None) -> Path:
    agent = normalize_agent(agent)
    if agent not in SUPPORTED_AGENTS:
        raise InstallError(f"unsupported agent: {agent}")
    if scope == "project":
        base = (project_dir or Path.cwd()).resolve()
        return base / (".agents/skills" if agent == "codex" else ".claude/skills")
    if scope != "user":
        raise InstallError(f"unsupported scope: {scope}")
    if agent == "codex":
        return Path.home() / ".agents" / "skills"
    else:
        config_root = Path(os.environ.get("CLAUDE_CONFIG_DIR", str(Path.home() / ".claude"))).expanduser()
    return config_root / "skills"


def state_path(scope: str, project_dir: Path | None = None) -> Path:
    if scope == "project":
        return (project_dir or Path.cwd()).resolve() / ".anywork" / "state.json"
    configured = os.environ.get("ANYWORK_HOME")
    return (Path(configured).expanduser() if configured else Path.home() / ".anywork") / "state.json"


def hash_tree(path: Path, agent: str | None = None) -> str:
    digest = hashlib.sha256()
    for child in sorted(item for item in path.rglob("*") if item.is_file()):
        relative = child.relative_to(path).as_posix()
        if "__pycache__" in child.parts or relative.endswith((".pyc", ".pyo")):
            continue
        if agent == "claude-code" and relative.startswith("agents/"):
            continue
        digest.update(relative.encode("utf-8"))
        digest.update(b"\0")
        digest.update(child.read_bytes())
        digest.update(b"\0")
    return digest.hexdigest()


def _load_state(path: Path) -> dict:
    if not path.exists():
        return {"schema_version": 1, "installations": {}}
    try:
        state = json.loads(path.read_text(encoding="utf-8"))
    except (json.JSONDecodeError, OSError) as error:
        raise InstallError(f"cannot read state file {path}: {error}") from error
    if state.get("schema_version") != 1 or not isinstance(state.get("installations"), dict):
        raise InstallError(f"invalid state file: {path}")
    return state


def _write_state(path: Path, state: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    payload = json.dumps(state, ensure_ascii=False, indent=2, sort_keys=True) + "\n"
    with tempfile.NamedTemporaryFile("w", encoding="utf-8", dir=path.parent, delete=False) as handle:
        handle.write(payload)
        temporary = Path(handle.name)
    temporary.replace(path)


def _backup_path(state_file: Path, agent: str, skill_id: str) -> Path:
    stamp = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%SZ")
    root = state_file.parent / "backups" / stamp / agent
    candidate = root / skill_id
    suffix = 1
    while candidate.exists():
        candidate = root / f"{skill_id}-{suffix}"
        suffix += 1
    return candidate


def _installation_key(agent: str, target: Path) -> str:
    return f"{agent}:{target.resolve()}"


def _materialize(source: Path, target: Path, agent: str) -> None:
    """Copy a canonical Skill and remove host-specific files for other agents."""
    shutil.copytree(source, target)
    if agent == "claude-code":
        codex_metadata = target / "agents"
        if codex_metadata.exists():
            shutil.rmtree(codex_metadata)


def install_skills(
    catalog: Catalog,
    skill_ids: list[str],
    agents: list[str],
    scope: str,
    project_dir: Path | None,
    dry_run: bool,
    force: bool,
    emit: Callable[[str, dict], None],
) -> None:
    state_file = state_path(scope, project_dir)
    state = _load_state(state_file)
    for raw_agent in agents:
        agent = normalize_agent(raw_agent)
        root = target_root(agent, scope, project_dir)
        for skill_id in skill_ids:
            source = catalog.skill_path(skill_id)
            target = root / skill_id
            source_hash = hash_tree(source, agent=agent)
            key = _installation_key(agent, target)
            record = state["installations"].get(key)
            if target.exists():
                current_hash = hash_tree(target)
                if current_hash == source_hash:
                    emit("unchanged", {"skill": skill_id, "agent": agent, "target": target})
                    if not dry_run and not record:
                        state["installations"][key] = {
                            "agent": agent,
                            "skill": skill_id,
                            "target": str(target),
                            "hash": source_hash,
                            "version": catalog.skills[skill_id]["version"],
                        }
                        _write_state(state_file, state)
                    continue
                owned_and_clean = bool(record and record.get("hash") == current_hash)
                if not owned_and_clean and not force:
                    raise InstallError(f"conflict:{target}")
                emit("plan_replace" if dry_run else "updated", {"skill": skill_id, "agent": agent, "target": target})
                if dry_run:
                    continue
            else:
                emit("plan_install" if dry_run else "installed", {"skill": skill_id, "agent": agent, "target": target})
                if dry_run:
                    continue
            root.mkdir(parents=True, exist_ok=True)
            temporary = Path(tempfile.mkdtemp(prefix=f".anywork-{skill_id}-", dir=root))
            backup: Path | None = None
            try:
                shutil.rmtree(temporary)
                _materialize(source, temporary, agent)
                if target.exists():
                    backup = _backup_path(state_file, agent, skill_id)
                    backup.parent.mkdir(parents=True, exist_ok=True)
                    shutil.move(str(target), str(backup))
                    emit("backup_note", {"backup": backup})
                temporary.replace(target)
            except Exception:
                if temporary.exists():
                    shutil.rmtree(temporary)
                if backup and backup.exists() and not target.exists():
                    shutil.move(str(backup), str(target))
                raise
            state["installations"][key] = {
                "agent": agent,
                "skill": skill_id,
                "target": str(target),
                "hash": source_hash,
                "version": catalog.skills[skill_id]["version"],
                "installed_at": datetime.now(timezone.utc).isoformat(),
            }
            _write_state(state_file, state)


def uninstall_skills(
    skill_ids: list[str],
    agents: list[str],
    scope: str,
    project_dir: Path | None,
    dry_run: bool,
    force: bool,
    emit: Callable[[str, dict], None],
) -> None:
    state_file = state_path(scope, project_dir)
    state = _load_state(state_file)
    for raw_agent in agents:
        agent = normalize_agent(raw_agent)
        root = target_root(agent, scope, project_dir)
        for skill_id in skill_ids:
            target = root / skill_id
            key = _installation_key(agent, target)
            record = state["installations"].get(key)
            if not record or not target.exists():
                emit("nothing_installed", {"skill": skill_id, "agent": agent, "target": target})
                continue
            current_hash = hash_tree(target)
            if current_hash != record.get("hash") and not force:
                raise InstallError(f"not_owned:{target}")
            backup = _backup_path(state_file, agent, skill_id)
            emit("plan_remove" if dry_run else "removed", {
                "skill": skill_id,
                "agent": agent,
                "target": target,
                "backup": backup,
            })
            if dry_run:
                continue
            backup.parent.mkdir(parents=True, exist_ok=True)
            shutil.move(str(target), str(backup))
            del state["installations"][key]
            _write_state(state_file, state)
