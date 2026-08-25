import logging
from sqlalchemy.orm import Session
from src.db import QuestPlanDB
from src.domain.plan_registry import PlanRegistry
from src.exceptions.generic import NotFoundError

logger = logging.getLogger(__name__)


class PlanRepository:
    def get_by_id(self, db: Session, plan_id: str) -> PlanRegistry | None:
        db_plan = db.get(QuestPlanDB, plan_id)
        if db_plan is None:
            return None
        return PlanRegistry.model_validate(db_plan)

    def get_all(self, db: Session) -> list[PlanRegistry]:
        db_plans = db.query(QuestPlanDB).all()
        return [PlanRegistry.model_validate(p) for p in db_plans]

    def add(self, db: Session, registry: PlanRegistry) -> PlanRegistry:
        db_plan = QuestPlanDB(**registry.model_dump())
        db.add(db_plan)
        db.flush()  # to get ID if auto-generated
        db.refresh(db_plan)
        logger.debug(f"Added plan '{db_plan.id}' to database")
        return PlanRegistry.model_validate(db_plan)

    def update(self, db: Session, registry: PlanRegistry) -> PlanRegistry:
        db_plan = db.get(QuestPlanDB, registry.id)
        if db_plan is None:
            raise NotFoundError("Plan", registry.id)
        for key, value in registry.model_dump().items():
            setattr(db_plan, key, value)
        db.commit()
        db.refresh(db_plan)
        return PlanRegistry.model_validate(db_plan)
