from __future__ import annotations

import os
import unittest
from unittest.mock import patch

from anywork.i18n import localized, message, normalize_language


class I18nTests(unittest.TestCase):
    def test_language_variants_normalize(self) -> None:
        self.assertEqual(normalize_language("zh_CN.UTF-8"), "zh-CN")
        self.assertEqual(normalize_language("ja-JP"), "ja")
        self.assertEqual(normalize_language("en-US"), "en")

    def test_environment_precedes_os_locale(self) -> None:
        with patch.dict(os.environ, {"ANYWORK_LOCALE": "ja"}):
            self.assertEqual(normalize_language(None), "ja")

    def test_messages_and_catalog_values_are_localized(self) -> None:
        self.assertIn("诊断", message("zh-CN", "doctor_fail", count=1))
        self.assertIn("人工", message("zh-CN", "quality_notice"))
        self.assertIn("実験版", message("ja", "quality_notice"))
        self.assertEqual(localized({"en": "Plan", "ja": "計画"}, "ja"), "計画")


if __name__ == "__main__":
    unittest.main()
