from typing import Any
from uuid import UUID
from sqlalchemy.orm import Session

from src.db.quest_settings import QuestSettingsDB

class QuestSettingsRepository:
    @staticmethod
    def add(db: Session, quest_id: UUID, form_data: dict[str, Any]) -> QuestSettingsDB:
        settings = QuestSettingsDB(quest_id=quest_id, form_data=form_data)
        db.add(settings)
        db.flush()
        db.refresh(settings)
        return settings

    @staticmethod
    def get_latest(db: Session, quest_id: UUID) -> QuestSettingsDB | None:
        return (
            db.query(QuestSettingsDB)
            .filter(QuestSettingsDB.quest_id == quest_id)
            .order_by(QuestSettingsDB.created_at.desc())
            .first()
        )
