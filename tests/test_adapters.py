from __future__ import annotations

import json
import unittest

from anywork.catalog import load_catalog


class AdapterMatrixTests(unittest.TestCase):
    def setUp(self) -> None:
        self.catalog = load_catalog()
        self.root = self.catalog.root
        self.matrix = json.loads((self.root / "registry" / "support-matrix.json").read_text(encoding="utf-8"))

    def test_stable_matrix_matches_catalog_and_descriptors(self) -> None:
        agents = {entry["id"]: entry for entry in self.matrix["agents"]}
        stable = {agent_id for agent_id, entry in agents.items() if entry["status"] == "stable"}
        self.assertEqual(stable, set(self.catalog.agents))
        self.assertEqual(agents["cursor"]["status"], "experimental")
        self.assertFalse(agents["cursor"]["default_install"])
        for agent_id, matrix_entry in agents.items():
            descriptor = json.loads((self.root / "adapters" / f"{agent_id}.json").read_text(encoding="utf-8"))
            self.assertEqual(descriptor["id"], agent_id)
            self.assertEqual(descriptor["status"], matrix_entry["status"])
            self.assertEqual(descriptor["target_family"], matrix_entry["target_family"])
            self.assertFalse(descriptor["install_executes_skill_code"])
            self.assertEqual(descriptor["skill_file"], "SKILL.md")
            self.assertTrue(descriptor["official_documentation"].startswith("https://"))

    def test_shared_agent_skills_family_has_one_physical_target(self) -> None:
        shared = [entry for entry in self.matrix["agents"] if entry["target_family"] == "agent-skills"]
        self.assertEqual({entry["project_target"] for entry in shared}, {".agents/skills"})
        self.assertEqual({entry["user_target"] for entry in shared}, {"~/.agents/skills"})

    def test_no_adapter_preapproves_tools(self) -> None:
        for path in (self.root / "adapters").glob("*.json"):
            descriptor = json.loads(path.read_text(encoding="utf-8"))
            self.assertNotIn("allowed_tools", descriptor)
            self.assertIn("portable_frontmatter_allowlist", descriptor)


if __name__ == "__main__":
    unittest.main()
