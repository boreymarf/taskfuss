import logging
from typing import Any
from sqlalchemy.orm import Session

from src.domain.quest_settings import QuestSettings
from src.repositories.quest_settings import QuestSettingsRepository

logger = logging.getLogger(__name__)

class QuestSettingsService:
    def __init__(self, quest_settings_repository: QuestSettingsRepository):
        self.quest_settings_repository = quest_settings_repository

    def add(self, db: Session, quest_id: int, form_data: dict[str, Any]) -> QuestSettings:
        settings = self.quest_settings_repository.add(db, quest_id, form_data)
        logger.debug(f"Added new quest settings for quest id '{quest_id}'")
        return settings

    def get_latest(self, db: Session, quest_id: int) -> QuestSettings | None:
        settings = self.quest_settings_repository.get_latest(db, quest_id)
        if settings is None:
            logger.debug(f"No settings found for quest id '{quest_id}'")
            return None
        return settings
