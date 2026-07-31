from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import (
    BAD_REQUEST,
    NOT_FOUND,
    UNAUTHORIZED,
)
from src.domain.quest import Quest, QuestCreateRequest
from src.service.quest import QuestService

router = APIRouter(prefix="/api/quest", tags=["quest"])


@router.post(
    "/",
    response_model=Quest,
    responses={**NOT_FOUND, **UNAUTHORIZED, **BAD_REQUEST},
)
def create_quest(
    data: QuestCreateRequest,
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    return QuestService.create_quest(db, current_user, data, auto_commit=True)

@router.get(
    "/",
    response_model=list[Quest],
    responses={**UNAUTHORIZED}
)
def get_all_quests(
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user)
) -> list[Quest]:
    return QuestService.get_all_by_user(db, current_user)


