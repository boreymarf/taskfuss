from pathlib import Path

import pytest
from sqlalchemy.orm import Session

from src.quests.registry import PlanRegistry
from src.service.quest_plan import QuestPlanService


@pytest.mark.usefixtures("db_session_class")
class TestQuestServiceWithDB:
    def test_add_and_get_plan(
        self, db_session_class: Session, project_root: Path
    ) -> None:
        assets = project_root / "tests" / "assets" / "plan_implementations"
        if not assets.exists():
            pytest.skip("Test plan directory not found")

        registry = PlanRegistry(
            id="test_plan",
            name="Test Plan",
            class_path="some.module:SomePlan",
        )
        db_plan = QuestPlanService.add_plan_to_db(db_session_class, registry)

        assert db_plan.id == "test_plan"
        fetched = QuestPlanService.get_plan(db_session_class, "test_plan")
        assert fetched is not None
        assert fetched.id == "test_plan"
        assert fetched.name == "Test Plan"

    def test_get_all_plans(self, db_session_class: Session) -> None:
        QuestPlanService.add_plan_to_db(
            db_session_class,
            PlanRegistry(
                id="p1",
                name="P1",
                class_path="a:b",
            ),
        )
        QuestPlanService.add_plan_to_db(
            db_session_class,
            PlanRegistry(
                id="p2",
                name="P2",
                class_path="c:d",
            ),
        )
        plans = QuestPlanService.get_all_plans(db_session_class)
        ids = {p.id for p in plans}
        assert ids >= {"p1", "p2"}

    def test_sync_plans_from_directory_creates_new(
        self, db_session_class: Session, project_root: Path
    ) -> None:
        assets = project_root / "tests" / "assets" / "plan_implementations"
        if not assets.exists():
            pytest.skip("Test plan directory not found")

        result = QuestPlanService.sync_plans_from_directory(
            db_session_class, assets, update_existing=False
        )
        assert "plan_a" in result
        assert "plan_b" in result
        db_plan = QuestPlanService.get_plan(db_session_class, "plan_a")
        assert db_plan is not None
        assert db_plan.name == "name_a"

    def test_sync_plans_updates_existing(
        self, db_session_class: Session, project_root: Path
    ) -> None:
        assets = project_root / "tests" / "assets" / "plan_implementations"
        if not assets.exists():
            pytest.skip("Test plan directory not found")

        QuestPlanService.sync_plans_from_directory(
            db_session_class, assets, update_existing=False
        )
        old = QuestPlanService.get_plan(db_session_class, "plan_a")
        assert old is not None

        old.name = "Modified"
        db_session_class.commit()

        result = QuestPlanService.sync_plans_from_directory(
            db_session_class, assets, update_existing=True
        )
        updated = result["plan_a"]
        assert updated is not None
        assert updated.name == "name_a"
        assert updated.name != "Modified"
