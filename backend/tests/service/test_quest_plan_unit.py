from pathlib import Path

import pytest

from src.quests.registry import PlanRegistry
from src.service.quest_plan import QuestPlanService
from tests.assets.plan_implementations.plan_a import PlanA


class TestLoadRegistriesFromFile:
    def test_loads_real_plan_file(self, project_root: Path):
        path = project_root / "tests" / "assets" / "plan_implementations" / "plan_a.py"
        if not path.exists():
            pytest.skip("Required plan file not found")

        registries = QuestPlanService.load_registries_from_file(path)

        assert len(registries) == 1
        registry = registries[0]
        assert isinstance(registry, PlanRegistry)
        assert registry.id
        assert registry.name

    def test_skips_init_file(self, project_root: Path):
        registries = QuestPlanService.load_registries_from_file(
            project_root / "tests" / "assets" / "plan_implementations" / "__init__.py"
        )
        assert registries == []

    def test_skips_non_py_file(self, tmp_path: Path):
        txt = tmp_path / "readme.txt"
        txt.write_text("nothing")
        registries = QuestPlanService.load_registries_from_file(txt)
        assert registries == []

    def test_empty_for_nonexistent_file(self, project_root: Path):
        path = project_root / "nonexistent.py"
        registries = QuestPlanService.load_registries_from_file(path)
        assert registries == []


class TestLoadRegistriesFromDirectory:
    def test_scans_plan_implementations(self, project_root: Path):
        assets = project_root / "tests" / "assets" / "plan_implementations"
        if not assets.exists():
            pytest.skip("Test plan directory not found")

        registries = QuestPlanService.load_registries_from_directory(assets)

        ids = {r.id for r in registries}
        assert len(ids) == 2
        assert "plan_a" in ids
