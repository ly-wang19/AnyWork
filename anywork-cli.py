#!/usr/bin/env python3
from __future__ import annotations

import sys
from pathlib import Path


REPOSITORY = Path(__file__).resolve().parent
sys.path.insert(0, str(REPOSITORY / "src"))

from anywork.cli import main  # noqa: E402


if __name__ == "__main__":
    raise SystemExit(main())
