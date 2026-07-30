import logging
from uuid import UUID

from sqlalchemy.orm import Session

from src.db import QuestDB

logger = logging.getLogger(__name__)


class QuestRepository:
    @staticmethod
    def add(db: Session, owner_id: int, plan_id: str) -> QuestDB:
        quest = QuestDB(owner_id=owner_id, plan_id=plan_id)
        db.add(quest)
        db.flush()
        db.refresh(quest)
        logger.debug(f"Added new quest with id '{quest.id}' for plan '{plan_id}'")
        return quest

    @staticmethod
    def get_by_id(db: Session, quest_id: UUID) -> QuestDB | None:
        return db.get(QuestDB, quest_id)
