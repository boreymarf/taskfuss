from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from src.api.dependencies.session import get_session
from src.exceptions import NotFoundError
from src.quests.fields import FieldForm
from src.quests.registry import PlanRegistryPublic
from src.service.quest_plan import QuestPlanService

router = APIRouter(prefix="/api/quest_plan", tags=["quest_plan"])


@router.get("/", response_model=list[PlanRegistryPublic])
def get_all_quest_plans(
    db: Session = Depends(get_session),
) -> list[PlanRegistryPublic]:
    plans = QuestPlanService.get_all_plans(db)
    return [plan.to_public() for plan in plans]


@router.get("/{plan_id}/setup_form", response_model=FieldForm)
def get_plan_setup_form(
    plan_id: str,
    db: Session = Depends(get_session),
) -> FieldForm:
    plan = QuestPlanService.get_plan(db, plan_id)

    if not plan:
        raise NotFoundError("Plan registry", plan_id)

    plan_class = plan.import_class()
    instance = plan_class()
    setup_form = instance.get_setup_form()
    return setup_form
