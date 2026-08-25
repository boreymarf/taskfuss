# src/services/plan_service.py

import importlib.util
import inspect
import logging
import sys
from pathlib import Path

from sqlalchemy.orm import Session

from src.domain.plan_registry import PlanRegistry
from src.exceptions import NotFoundError
from src.plans.base import BasePlan
from src.repositories.plan import PlanRepository

logger = logging.getLogger(__name__)


class PlanService:
    def __init__(self, plan_repository: PlanRepository):
        self.plan_repository = plan_repository

    def load_registries_from_directory(self, directory: Path) -> list[PlanRegistry]:
        """Load plan registries from all Python files in a directory recursively."""
        registries: list[PlanRegistry] = []
        for py_file in directory.rglob("*.py"):
            registries.extend(self.load_registries_from_file(py_file))
        return registries

    def load_registries_from_file(self, path: Path) -> list[PlanRegistry]:
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

    def sync_plans_from_directory(
        self,
        db: Session,
        directory: Path,
        *,
        update_existing: bool = True,
    ) -> dict[str, PlanRegistry]:
        """Synchronize plans from directory to database, adding new and optionally updating existing."""
        registries = self.load_registries_from_directory(directory)
        result: dict[str, PlanRegistry] = {}

        for registry in registries:
            existing = self.plan_repository.get_by_id(db, registry.id)
            if existing is None:
                db_plan = self.plan_repository.add(db, registry)
            elif update_existing:
                db_plan = self.plan_repository.update(db, registry)
            else:
                db_plan = existing
            result[registry.id] = db_plan

        return result

    def get_plan(self, db: Session, plan_id: str) -> PlanRegistry | None:
        """Retrieve a single plan by ID."""
        return self.plan_repository.get_by_id(db, plan_id)

    def get_all_plans(self, db: Session) -> list[PlanRegistry]:
        """Retrieve all plans from the database."""
        return self.plan_repository.get_all(db)

    def get_plan_instance(self, db: Session, plan_id: str) -> BasePlan:
        """Retrieve a plan from the database and return an instantiated plan object."""
        plan_registry = self.get_plan(db, plan_id)
        if plan_registry is None:
            raise NotFoundError("Plan", plan_id)

        plan_cls = plan_registry.import_class()
        return plan_cls()
