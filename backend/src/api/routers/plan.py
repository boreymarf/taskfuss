from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from src.api.dependencies.session import get_session
from src.api.openapi_responses import NOT_FOUND
from src.domain.fields import FormFields
from src.domain.plan_registry import PlanRegistryPublic
from src.exceptions import NotFoundError
from src.service.plan import PlanService

router = APIRouter(prefix="/api/plan", tags=["plan"])


@router.get("/", response_model=list[PlanRegistryPublic])
def get_all_plans(
    db: Session = Depends(get_session),
) -> list[PlanRegistryPublic]:
    plans = PlanService.get_all_plans(db)
    return [plan.to_public() for plan in plans]


@router.get(
    "/{plan_id}/settings_form", response_model=FormFields, responses={**NOT_FOUND}
)
def get_plan_settings_form(
    plan_id: str,
    db: Session = Depends(get_session),
) -> FormFields:
    plan = PlanService.get_plan(db, plan_id)

    if not plan:
        raise NotFoundError("Plan registry", plan_id)

    plan_class = plan.import_class()
    instance = plan_class()
    setup_form = instance.get_settings_form()
    return setup_form
