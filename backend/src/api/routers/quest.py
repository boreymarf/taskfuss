from typing import Annotated

from dependency_injector.wiring import Provide, inject
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session
from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import BAD_REQUEST, NOT_FOUND, UNAUTHORIZED
from src.container import Container
from src.domain.quest import Quest, QuestCreateRequest
from src.service.quest import QuestService

router = APIRouter(prefix="/api/quest", tags=["quest"])

@router.post("/", response_model=Quest, responses={**NOT_FOUND, **UNAUTHORIZED, **BAD_REQUEST})
@inject
def create_quest(
    data: QuestCreateRequest,
    quest_service: Annotated[QuestService, Depends(Provide[Container.quest_service])],
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return quest_service.create_quest(db, current_user, data, auto_commit=True)


@router.get("/", response_model=list[Quest], responses={**UNAUTHORIZED})
@inject
def get_all_quests(
    quest_service: Annotated[QuestService, Depends(Provide[Container.quest_service])],
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return quest_service.get_all_by_user(db, current_user)
