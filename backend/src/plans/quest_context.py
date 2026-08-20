from datetime import datetime
import logging

from sqlalchemy.orm import Session

from src.domain.quest_settings import QuestSettings
from src.domain.quest_state import QuestState
from src.domain.record import Record
from src.repositories.quest_settings import QuestSettingsRepository
from src.service.quest_state import QuestStateService
from src.service.record import RecordService

logger = logging.getLogger(__name__)


class QuestContext:

    def __init__(
        self,
        db: Session,
        quest_id: int,
    ):
        self.db = db
        self.quest_id = quest_id

    def get_quest_id(self) -> int:
        return self.quest_id

    def get_quest_settings(self) -> QuestSettings | None:
        quest_settings_db = QuestSettingsRepository.get_latest(self.db, self.quest_id)

        if not quest_settings_db:
            return None

        quest_settings = QuestSettings.model_validate(quest_settings_db)
        return quest_settings

    def get_state_for_date(self, date: datetime) -> QuestState | None:
        return QuestStateService.get_by_date(self.db, self.quest_id, date)

    # TODO
    def get_latest_records_for_state(self, state_id: int, *, include_defaults: bool =True) -> list[Record]:

        records = RecordService.get_all(
            self.db, quest_id=self.quest_id, state_id=state_id, latest=True
        )

        return records
