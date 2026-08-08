from __future__ import annotations

import io
import tempfile
import unittest
from contextlib import redirect_stdout
from pathlib import Path

from anywork.cli import main
from anywork.experience import DEMOS, render_demo


class ExperienceTests(unittest.TestCase):
    def test_demos_are_trilingual_and_use_the_orchestrator(self) -> None:
        self.assertEqual({demo.id for demo in DEMOS}, {"research", "meeting", "automation"})
        for demo in DEMOS:
            self.assertEqual(set(demo.name), {"en", "zh-CN", "ja"})
            self.assertEqual(set(demo.prompt), {"en", "zh-CN", "ja"})
            self.assertEqual(set(demo.outcome), {"en", "zh-CN", "ja"})
            self.assertEqual(demo.chain[0], "orchestrate-work")
            for language in ("en", "zh-CN", "ja"):
                self.assertIn("$orchestrate-work", demo.prompt[language])
                self.assertIn("orchestrate-work", render_demo(language, demo.id))

    def test_setup_dry_run_installs_nothing(self) -> None:
        with tempfile.TemporaryDirectory() as temporary:
            project = Path(temporary)
            output = io.StringIO()
            with redirect_stdout(output):
                result = main(
                    [
                        "setup",
                        "--scope",
                        "project",
                        "--project-dir",
                        str(project),
                        "--agent",
                        "codex",
                        "--dry-run",
                        "--lang",
                        "en",
                    ]
                )
            self.assertEqual(result, 0)
            self.assertIn("orchestrate-work", output.getvalue())
            self.assertFalse((project / ".agents").exists())


if __name__ == "__main__":
    unittest.main()
