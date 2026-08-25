import logging
from typing import Any
from sqlalchemy.orm import Session

from src.db.quest_settings import QuestSettingsDB
from src.domain.quest_settings import QuestSettings

logger = logging.getLogger(__name__)

class QuestSettingsRepository:
    def add(self, db: Session, quest_id: int, form_data: dict[str, Any]) -> QuestSettings:
        settings_db = QuestSettingsDB(quest_id=quest_id, form_data=form_data)
        db.add(settings_db)
        db.flush()
        db.refresh(settings_db)
        logger.debug(f"Added new quest settings for quest id '{quest_id}'")
        return QuestSettings.model_validate(settings_db)

    def get_latest(self, db: Session, quest_id: int) -> QuestSettings | None:
        settings_db = (
            db.query(QuestSettingsDB)
            .filter(QuestSettingsDB.quest_id == quest_id)
            .order_by(QuestSettingsDB.created_at.desc())
            .first()
        )
        if settings_db is None:
            return None
        return QuestSettings.model_validate(settings_db)
