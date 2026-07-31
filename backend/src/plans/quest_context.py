from datetime import date, datetime
from uuid import UUID

from sqlalchemy.orm import Session

from src.domain.quest_settings import QuestSettings
from src.domain.record import Record
from src.repositories.quest_settings import QuestSettingsRepository


class QuestContext:

    def __init__(
        self,
        session: Session,
        quest_id: UUID,
    ):
        self.session = session
        self.quest_id = quest_id

    def get_quest_id(self) -> UUID:
        return self.quest_id

    def get_quest_settings(self) -> QuestSettings | None:
        quest_settings_db = QuestSettingsRepository.get_latest(
            self.session, self.quest_id
        )

        if not quest_settings_db:
            return None

        quest_settings = QuestSettings.model_validate(quest_settings_db)
        return quest_settings
