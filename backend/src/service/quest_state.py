from datetime import datetime

from sqlalchemy.orm import Session

from src.domain.quest_state import QuestState, QuestStateCreate
from src.repositories.quest_state import QuestStateRepository

import logging

logger = logging.getLogger(__name__)


class QuestStateService:
    @staticmethod
    def get_all_current(db: Session, owner_id: int) -> list[QuestState]:
        states_db = QuestStateRepository.get_all(
            db,
            owner_id,
            start_date=datetime.now(),
        )
        logger.debug(
            f"Retrieved {len(states_db)} current states for owner_id={owner_id}"
        )
        return [QuestState.model_validate(state) for state in states_db]

    @staticmethod
    def get_all(db: Session, owner_id: int) -> list[QuestState]:
        states_db = QuestStateRepository.get_all(db, owner_id)
        logger.debug(f"Retrieved {len(states_db)} states for owner_id={owner_id}")
        return [QuestState.model_validate(state) for state in states_db]

    @staticmethod
    def add(db: Session, state_data: QuestStateCreate) -> QuestState:
        state_db = QuestStateRepository.add(db, state_data)
        logger.debug(f"Added quest state for quest '{state_data.quest_id}'")
        return QuestState.model_validate(state_db)

    @staticmethod
    def get_by_id(db: Session, state_id: int) -> QuestState | None:
        state_db = QuestStateRepository.get_by_id(db, state_id)
        if state_db is None:
            logger.debug(f"Quest state with id {state_id} not found")
            return None
        return QuestState.model_validate(state_db)

    @staticmethod
    def get_latest(db: Session, quest_id: int) -> QuestState | None:
        state_db = QuestStateRepository.get_latest(db, quest_id)
        if state_db is None:
            logger.debug(f"No latest state found for quest {quest_id}")
            return None
        return QuestState.model_validate(state_db)

    @staticmethod
    def has_overlap(
        db: Session,
        quest_id: int,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
    ) -> bool:
        result = QuestStateRepository.has_overlap(db, quest_id, start_date, end_date)
        logger.debug(
            f"Overlap check for quest {quest_id} between {start_date} and {end_date}: {result}"
        )
        return result

    @staticmethod
    def get_by_date(db: Session, quest_id: int, target_date: datetime) -> QuestState | None:
        state_db = QuestStateRepository.get_by_date(db, quest_id, target_date)
        if state_db is None:
            logger.debug(f"No state found for quest {quest_id} at {target_date}")
            return None
        return QuestState.model_validate(state_db)
