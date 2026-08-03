import logging
from typing import Any
from sqlalchemy.orm import Session

from src.db.quest_settings import QuestSettingsDB

logger = logging.getLogger(__name__)

class QuestSettingsRepository:
    @staticmethod
    def add(db: Session, quest_id: int, form_data: dict[str, Any]) -> QuestSettingsDB:
        settings = QuestSettingsDB(quest_id=quest_id, form_data=form_data)
        db.add(settings)
        db.flush()
        db.refresh(settings)
        logger.debug(f"Added new quest settings for quest id '{quest_id}'")
        return settings

    @staticmethod
    def get_latest(db: Session, quest_id: int) -> QuestSettingsDB | None:
        return (
            db.query(QuestSettingsDB)
            .filter(QuestSettingsDB.quest_id == quest_id)
            .order_by(QuestSettingsDB.created_at.desc())
            .first()
        )
