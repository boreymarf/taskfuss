from typing import Annotated

from dependency_injector.wiring import Provide, inject
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import UNAUTHORIZED
from src.container import Container
from src.domain.quest_state import QuestState
from src.service.quest_state import QuestStateService

router = APIRouter(prefix="/api/quest_state", tags=["quest_state"])

@router.get("/today", response_model=list[QuestState], responses={**UNAUTHORIZED})
@inject
def get_states_today(
    state_service: Annotated[QuestStateService, Depends(Provide[Container.quest_state_service])],
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return state_service.get_all_current(db, current_user)


@router.get("/", response_model=list[QuestState], responses={**UNAUTHORIZED})
@inject
def get_all_states(
    state_service: Annotated[QuestStateService, Depends(Provide[Container.quest_state_service])],
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return state_service.get_all(db, current_user)
