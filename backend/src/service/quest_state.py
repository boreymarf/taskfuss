# src/services/quest_state_service.py

from datetime import datetime
from sqlalchemy.orm import Session

from src.domain.quest_state import QuestState, QuestStateCreate
from src.repositories.quest_state import QuestStateRepository

import logging

logger = logging.getLogger(__name__)

class QuestStateService:
    def __init__(self, quest_state_repository: QuestStateRepository):
        self.quest_state_repository = quest_state_repository

    def get_all_current(self, db: Session, owner_id: int) -> list[QuestState]:
        states = self.quest_state_repository.get_all(
            db,
            owner_id,
            start_date=datetime.now(),
        )
        logger.debug(
            f"Retrieved {len(states)} current states for owner_id={owner_id}"
        )
        return states

    def get_all(self, db: Session, owner_id: int) -> list[QuestState]:
        states = self.quest_state_repository.get_all(db, owner_id)
        logger.debug(f"Retrieved {len(states)} states for owner_id={owner_id}")
        return states

    def add(self, db: Session, state_data: QuestStateCreate) -> QuestState:
        state = self.quest_state_repository.add(db, state_data)
        logger.debug(f"Added quest state for quest '{state_data.quest_id}'")
        return state

    def get_by_id(self, db: Session, state_id: int) -> QuestState | None:
        state = self.quest_state_repository.get_by_id(db, state_id)
        if state is None:
            logger.debug(f"Quest state with id {state_id} not found")
            return None
        return state

    def get_latest(self, db: Session, quest_id: int) -> QuestState | None:
        state = self.quest_state_repository.get_latest(db, quest_id)
        if state is None:
            logger.debug(f"No latest state found for quest {quest_id}")
            return None
        return state

    def has_overlap(
        self,
        db: Session,
        quest_id: int,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
    ) -> bool:
        result = self.quest_state_repository.has_overlap(
            db, quest_id, start_date, end_date
        )
        logger.debug(
            f"Overlap check for quest {quest_id} between {start_date} and {end_date}: {result}"
        )
        return result

    def get_by_date(self, db: Session, quest_id: int, target_date: datetime) -> QuestState | None:
        state = self.quest_state_repository.get_by_date(db, quest_id, target_date)
        if state is None:
            logger.debug(f"No state found for quest {quest_id} at {target_date}")
            return None
        return state
