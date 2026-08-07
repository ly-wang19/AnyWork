from __future__ import annotations

import argparse
import sys
from pathlib import Path

from .catalog import load_catalog, validate_catalog
from .i18n import help_text, localized, message, normalize_language
from .installer import InstallError, SUPPORTED_AGENTS, install_skills, normalize_agent, uninstall_skills


def _extract_language(argv: list[str]) -> tuple[str, list[str]]:
    language: str | None = None
    cleaned: list[str] = []
    index = 0
    while index < len(argv):
        item = argv[index]
        if item == "--lang":
            if index + 1 >= len(argv):
                raise SystemExit("--lang requires a value")
            language = argv[index + 1]
            index += 2
            continue
        if item.startswith("--lang="):
            language = item.split("=", 1)[1]
            index += 1
            continue
        cleaned.append(item)
        index += 1
    if language:
        prefix = language.replace("_", "-").lower().split("-", 1)[0]
        if prefix not in {"en", "zh", "ja"}:
            raise SystemExit(f"unsupported language: {language}")
    normalized = normalize_language(language)
    return normalized, cleaned


def _add_help(parser: argparse.ArgumentParser, language: str) -> None:
    parser.add_argument("-h", "--help", action="help", help=help_text(language, "help"))


def build_parser(language: str) -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(
        prog="anywork",
        description=f"Any work. Any agent. Any machine. {help_text(language, 'description')}",
        add_help=False,
    )
    _add_help(parser, language)
    parser.add_argument("--version", action="version", version="anywork 0.1.0")
    subparsers = parser.add_subparsers(dest="command", required=True)

    list_parser = subparsers.add_parser("list", help=help_text(language, "list"), add_help=False)
    _add_help(list_parser, language)
    doctor_parser = subparsers.add_parser("doctor", help=help_text(language, "doctor"), add_help=False)
    _add_help(doctor_parser, language)

    for command in ("install", "update", "uninstall"):
        subparser = subparsers.add_parser(command, help=help_text(language, command), add_help=False)
        _add_help(subparser, language)
        subparser.add_argument("pack", help=help_text(language, "pack"))
        subparser.add_argument(
            "--agent",
            action="append",
            default=[],
            help=help_text(language, "agent"),
        )
        subparser.add_argument("--scope", choices=("user", "project"), default="user", help=help_text(language, "scope"))
        subparser.add_argument("--project-dir", type=Path, default=Path.cwd(), help=help_text(language, "project_dir"))
        subparser.add_argument("--dry-run", action="store_true", help=help_text(language, "dry_run"))
        subparser.add_argument("--force", action="store_true", help=help_text(language, "force"))
    return parser


def _selected_agents(values: list[str]) -> list[str]:
    if not values or "all" in values:
        return list(SUPPORTED_AGENTS)
    result: list[str] = []
    for value in values:
        agent = normalize_agent(value)
        if agent not in SUPPORTED_AGENTS:
            raise InstallError(f"agent_unknown:{value}")
        if agent not in result:
            result.append(agent)
    return result


def main(argv: list[str] | None = None) -> int:
    language, cleaned = _extract_language(list(argv if argv is not None else sys.argv[1:]))
    args = build_parser(language).parse_args(cleaned)
    catalog = load_catalog()

    def emit(key: str, values: dict) -> None:
        print(message(language, key, **values))

    if args.command == "list":
        print(message(language, "catalog_title"))
        for pack in catalog.raw.get("packs", []):
            print(message(
                language,
                "pack_line",
                name=localized(pack["display_name"], language),
                description=localized(pack["description"], language),
                count=len(catalog.resolve_pack(pack["id"])),
            ))
        return 0

    if args.command == "doctor":
        problems = validate_catalog(catalog)
        if problems:
            print(message(language, "doctor_fail", count=len(problems)))
            for problem in problems:
                print(message(language, "validation_error", message=problem))
            return 1
        print(message(
            language,
            "doctor_ok",
            skills=len(catalog.skills),
            packs=len(catalog.packs),
            agents=len(catalog.agents),
        ))
        return 0

    if args.pack not in catalog.packs:
        print(message(language, "pack_unknown", pack=args.pack), file=sys.stderr)
        return 2
    skill_ids = catalog.resolve_pack(args.pack)
    try:
        agents = _selected_agents(args.agent)
        if args.command in {"install", "update"}:
            install_skills(
                catalog=catalog,
                skill_ids=skill_ids,
                agents=agents,
                scope=args.scope,
                project_dir=args.project_dir,
                dry_run=args.dry_run,
                force=args.force,
                emit=emit,
            )
        else:
            uninstall_skills(
                skill_ids=skill_ids,
                agents=agents,
                scope=args.scope,
                project_dir=args.project_dir,
                dry_run=args.dry_run,
                force=args.force,
                emit=emit,
            )
    except InstallError as error:
        raw = str(error)
        if ":" in raw:
            key, value = raw.split(":", 1)
            field = "target" if key in {"conflict", "not_owned"} else "agent"
            print(message(language, key, **{field: value}), file=sys.stderr)
        else:
            print(raw, file=sys.stderr)
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
