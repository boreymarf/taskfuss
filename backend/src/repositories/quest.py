import logging

from sqlalchemy.orm import Session

from src.db import QuestDB
from src.domain.quest import Quest

logger = logging.getLogger(__name__)


class QuestRepository:
    @staticmethod
    def add(db: Session, owner_id: int, plan_id: str) -> Quest:
        quest_db = QuestDB(owner_id=owner_id, plan_id=plan_id)
        db.add(quest_db)
        db.flush()
        db.refresh(quest_db)
        logger.debug(f"Added new quest with id '{quest_db.id}' for plan '{plan_id}'")
        return Quest.model_validate(quest_db)

    @staticmethod
    def get_by_id(db: Session, quest_id: int) -> Quest | None:
        quest_db = db.get(QuestDB, quest_id)
        return Quest.model_validate(quest_db) if quest_db else None

    @staticmethod
    def get_all_by_user(db: Session, owner_id: int) -> list[Quest]:
        quests_db = db.query(QuestDB).filter(QuestDB.owner_id == owner_id).all()
        return [Quest.model_validate(q) for q in quests_db]
