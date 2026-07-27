from datetime import date, datetime
from uuid import UUID

from pydantic import BaseModel
from sqlalchemy.orm import Session

from src.domain.record import Record
from src.plans.setup_data import QuestSettings


class QuestContext:

    def __init__(
        self,
        session: Session,
        quest_id: UUID,
    ):
        self._session = session
        self.quest_id = quest_id

    def get_setup_data(
        self
    ):
        pass

    def get_all_latest_records(
        self, field_path: str, date: datetime
    ) -> dict[str, Record]:
        return {}

    def get_all_records_for_date(
        self,
        date: date,
        limit: int = 100,
    ) -> dict[str, Record]:
        return {}
