
from datetime import datetime
from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import UNAUTHORIZED
from src.domain.quest_state import QuestState
from src.service.quest_state import QuestStateService

router = APIRouter(prefix="/api/quest_state", tags=["quest_state"])

@router.get(
    "/today",
    response_model=list[QuestState],
    responses={**UNAUTHORIZED}
)
def get_states_today(
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user)
) -> list[QuestState]:
    return QuestStateService.get_all_current(db, current_user)

@router.get(
    "/",
    response_model=list[QuestState],
    responses={**UNAUTHORIZED}
)
def get_all_states(
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user)
) -> list[QuestState]:
    return QuestStateService.get_all(db, current_user)
