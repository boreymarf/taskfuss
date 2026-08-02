from datetime import datetime
import logging
from uuid import UUID
from sqlalchemy.orm import Session

from src.db.record import RecordDB
from src.domain import Record
from src.domain.record import RecordCreateRequest
from src.exceptions.generic import ForbiddenError, NotFoundError
from src.repositories.quest import QuestRepository
from src.repositories.quest_state import QuestStateRepository
from src.repositories.record import RecordRepository

logger = logging.getLogger(__name__)


class RecordService:
    @staticmethod
    def get_all(db: Session) -> list[Record]:
        records_db = RecordRepository.get_all(db)
        result = [Record.model_validate(r) for r in records_db]
        logger.debug(f"Fetched {len(result)} records")
        return result

    @staticmethod
    def get_all_by_quest_id(db: Session, quest_id: UUID) -> list[Record]:
        records_db = RecordRepository.get_all_by_quest_id(db, quest_id)
        result = [Record.model_validate(r) for r in records_db]
        logger.debug(f"Fetched {len(result)} records for quest_id={quest_id}")
        return result

    @staticmethod
    def get(db: Session, record_id: UUID) -> Record | None:
        record_db = RecordRepository.get(db, record_id)
        if record_db is None:
            raise NotFoundError("Record", record_id)
        result = Record.model_validate(record_db)
        logger.debug(f"Record {record_id} found")
        return result

    @staticmethod
    def create(db: Session, request: RecordCreateRequest, user_id: int) -> Record:
        quest_db = QuestRepository.get_by_id(db, request.quest_id)

        if not quest_db:
            raise NotFoundError("Quest", request.quest_id)

        if quest_db.owner_id != user_id:
            raise ForbiddenError()

        # TODO: Replace later with actual record time
        quest_state_db = QuestStateRepository.get_by_date(
            db, request.quest_id, datetime.now()
        )

        if quest_state_db:
            raise NotFoundError("Quest state")



        record_db = RecordDB(**request.model_dump())
        created = RecordRepository.add(db, record_db)
        result = Record.model_validate(created)
        logger.debug(f"Created record with id={result.id}")
        return result

    @staticmethod
    def remove(db: Session, record_id: UUID) -> None:
        record_db = RecordRepository.get(db, record_id)
        if record_db is None:
            logger.debug(f"Record {record_id} not found, nothing to remove")
            return
        RecordRepository.remove(db, record_db)
        logger.debug(f"Removed record {record_id}")
