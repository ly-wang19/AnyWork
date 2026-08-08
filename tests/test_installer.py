from __future__ import annotations

import json
import tempfile
import unittest
from pathlib import Path

from anywork.catalog import load_catalog
from anywork.installer import InstallError, install_skills, uninstall_skills


class InstallerTests(unittest.TestCase):
    def setUp(self) -> None:
        self.catalog = load_catalog()
        self.temporary = tempfile.TemporaryDirectory()
        self.project = Path(self.temporary.name)
        self.events: list[tuple[str, dict]] = []

    def tearDown(self) -> None:
        self.temporary.cleanup()

    def emit(self, key: str, values: dict) -> None:
        self.events.append((key, values))

    def install(self, agents: list[str] | None = None, dry_run: bool = False) -> None:
        install_skills(
            catalog=self.catalog,
            skill_ids=["clarify-outcome"],
            agents=agents or ["codex", "claude-code"],
            scope="project",
            project_dir=self.project,
            dry_run=dry_run,
            force=False,
            emit=self.emit,
        )

    def test_dry_run_writes_nothing(self) -> None:
        self.install(dry_run=True)
        self.assertFalse((self.project / ".agents").exists())
        self.assertFalse((self.project / ".claude").exists())
        self.assertFalse((self.project / ".anywork").exists())

    def test_installs_for_both_agents_and_is_idempotent(self) -> None:
        self.install()
        self.assertTrue((self.project / ".agents/skills/clarify-outcome/SKILL.md").is_file())
        self.assertTrue((self.project / ".claude/skills/clarify-outcome/SKILL.md").is_file())
        self.assertTrue((self.project / ".agents/skills/clarify-outcome/agents/openai.yaml").is_file())
        self.assertFalse((self.project / ".claude/skills/clarify-outcome/agents").exists())
        state = json.loads((self.project / ".anywork/state.json").read_text(encoding="utf-8"))
        self.assertEqual(len(state["installations"]), 2)
        self.events.clear()
        self.install()
        self.assertEqual([event[0] for event in self.events], ["unchanged", "unchanged"])

    def test_shared_target_tracks_consumers_without_duplicate_files(self) -> None:
        shared_agents = ["codex", "gemini-cli", "github-copilot", "opencode"]
        self.install(agents=shared_agents)
        state = json.loads((self.project / ".anywork/state.json").read_text(encoding="utf-8"))
        self.assertEqual(len(state["installations"]), 1)
        record = next(iter(state["installations"].values()))
        self.assertEqual(set(record["consumers"]), set(shared_agents))
        uninstall_skills(
            skill_ids=["clarify-outcome"],
            agents=["gemini-cli"],
            scope="project",
            project_dir=self.project,
            dry_run=False,
            force=False,
            emit=self.emit,
        )
        self.assertTrue((self.project / ".agents/skills/clarify-outcome/SKILL.md").is_file())
        state = json.loads((self.project / ".anywork/state.json").read_text(encoding="utf-8"))
        record = next(iter(state["installations"].values()))
        self.assertNotIn("gemini-cli", record["consumers"])
        self.assertIn("codex", record["consumers"])

    def test_refuses_to_overwrite_local_changes(self) -> None:
        self.install(agents=["codex"])
        skill = self.project / ".agents/skills/clarify-outcome/SKILL.md"
        skill.write_text(skill.read_text(encoding="utf-8") + "\nlocal edit\n", encoding="utf-8")
        with self.assertRaisesRegex(InstallError, "conflict"):
            self.install(agents=["codex"])

    def test_uninstall_moves_managed_skill_to_backup(self) -> None:
        self.install(agents=["codex"])
        uninstall_skills(
            skill_ids=["clarify-outcome"],
            agents=["codex"],
            scope="project",
            project_dir=self.project,
            dry_run=False,
            force=False,
            emit=self.emit,
        )
        self.assertFalse((self.project / ".agents/skills/clarify-outcome").exists())
        backups = list((self.project / ".anywork/backups").rglob("clarify-outcome/SKILL.md"))
        self.assertEqual(len(backups), 1)


if __name__ == "__main__":
    unittest.main()
