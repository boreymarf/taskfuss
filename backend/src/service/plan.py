import importlib.util
import inspect
import logging
import sys
from pathlib import Path

from sqlalchemy.orm import Session

from src.db import QuestPlanDB
from src.domain.plan_registry import PlanRegistry
from src.exceptions import NotFoundError
from src.plans.base import BasePlan

logger = logging.getLogger(__name__)


class PlanService:
    @staticmethod
    def load_registries_from_directory(directory: Path) -> list[PlanRegistry]:
        """Load plan registries from all Python files in a directory recursively."""
        registries: list[PlanRegistry] = []
        for py_file in directory.rglob("*.py"):
            registries.extend(PlanService.load_registries_from_file(py_file))
        return registries

    @staticmethod
    def load_registries_from_file(path: Path) -> list[PlanRegistry]:
        """Load plan registries from a single Python file."""
        if path.suffix != ".py" or path.stem == "__init__":
            return []

        name = path.stem
        try:
            spec = importlib.util.spec_from_file_location(name, path)
            if spec is None or spec.loader is None:
                logger.warning("Cannot create spec for %s", path)
                return []
            ns = importlib.util.module_from_spec(spec)
            sys.modules[name] = ns
            spec.loader.exec_module(ns)
        except Exception as e:
            logger.warning("Failed to load file %s: %s", path, e)
            return []

        registries: list[PlanRegistry] = []
        for obj in vars(ns).values():
            if (
                inspect.isclass(obj)
                and issubclass(obj, BasePlan)
                and obj is not BasePlan
            ):
                instance = obj()
                meta = instance.get_metadata()
                registry = PlanRegistry(
                    id=meta.id,
                    name=meta.name,
                    class_path=f"{path}:{obj.__name__}",
                )
                registries.append(registry)

        if len(registries) > 0:
            logger.debug(
                f"Loaded {len(registries)} plan registries from '{path.stem}' file"
            )
        return registries

    @staticmethod
    def sync_plans_from_directory(
        db: Session,
        directory: Path,
        *,
        update_existing: bool = True,
    ) -> dict[str, QuestPlanDB]:
        """Synchronize plans from directory to database, adding new and optionally updating existing."""
        registries = PlanService.load_registries_from_directory(directory)
        result: dict[str, QuestPlanDB] = {}

        for registry in registries:
            existing = db.get(QuestPlanDB, registry.id)
            if existing is None:
                db_plan = PlanService.add_plan_to_db(db, registry)
            elif update_existing:
                for key, value in registry.model_dump().items():
                    setattr(existing, key, value)
                db.commit()
                db.refresh(existing)
                db_plan = existing
            else:
                db_plan = existing
            result[registry.id] = db_plan

        return result

    @staticmethod
    def get_plan(db: Session, plan_id: str) -> PlanRegistry | None:
        """Retrieve a single plan by ID."""
        row = db.get(QuestPlanDB, plan_id)
        if row is None:
            return None
        return PlanRegistry.model_validate(row)

    @staticmethod
    def get_all_plans(db: Session) -> list[PlanRegistry]:
        """Retrieve all plans from the database."""
        rows = db.query(QuestPlanDB).all()
        return [PlanRegistry.model_validate(row) for row in rows]

    @staticmethod
    def get_plan_instance(db: Session, plan_id: str) -> BasePlan:
        """Retrieve a plan from the database and return an instantiated plan object."""
        plan_registry = PlanService.get_plan(db, plan_id)
        if plan_registry is None:
            raise NotFoundError("Plan", plan_id)

        plan_cls = plan_registry.import_class()
        return plan_cls()

    @staticmethod
    def add_plan_to_db(db: Session, registry: PlanRegistry) -> QuestPlanDB:
        """Add a new plan registry to the database."""
        db_plan = QuestPlanDB(**registry.model_dump())
        db.add(db_plan)
        db.commit()
        db.refresh(db_plan)
        return db_plan
