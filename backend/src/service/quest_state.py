from datetime import datetime

from sqlalchemy.orm import Session

from src.domain.quest_state import QuestState
from src.repositories.quest_state import QuestStateRepository

import logging

logger = logging.getLogger(__name__)

class QuestStateService:
    @staticmethod
    def get_all_current(db: Session, owner_id: int) -> list[QuestState]:
        states_db = QuestStateRepository.get_all_by_date(db, datetime.now(), owner_id)
        logger.debug(f"Retrieved {len(states_db)} current states for owner_id={owner_id}")
        return [QuestState.model_validate(state) for state in states_db]
