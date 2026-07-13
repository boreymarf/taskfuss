from pprint import pprint

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from src.api.dependencies.auth import get_current_user
from src.api.dependencies.session import get_session
from src.api.openapi_responses import NOT_FOUND, PLAN_SETUP_VALIDATION_ERROR, UNAUTHORIZED
from src.dto.quest import QuestCreate
from src.service.quest import QuestService

router = APIRouter(prefix="/api/quest", tags=["quest"])


@router.post("/", responses={**NOT_FOUND, **UNAUTHORIZED, **PLAN_SETUP_VALIDATION_ERROR})
def create_quest(
    data: QuestCreate,
    db: Session = Depends(get_session),
    current_user: int = Depends(get_current_user),
):
    pprint(data)

    quest = QuestService.create_quest(db, current_user, data)
    db.commit()
    return quest.__dict__
