# src/services/record_service.py

from datetime import datetime
import logging
from sqlalchemy.orm import Session

from src.classes.form_processor import FormProcessor
from src.domain.record import Record, RecordCreate
from src.exceptions.generic import BadRequestError, ForbiddenError, NotFoundError
from src.exceptions.quest_state import NoFieldsQuestStateError
from src.exceptions.record import RecordValidationFailed
from src.repositories.quest import QuestRepository
from src.repositories.quest_state import QuestStateRepository
from src.repositories.record import RecordRepository

logger = logging.getLogger(__name__)

class RecordService:
    def __init__(
        self,
        record_repository: RecordRepository,
        quest_repository: QuestRepository,
        quest_state_repository: QuestStateRepository,
    ):
        self.record_repository = record_repository
        self.quest_repository = quest_repository
        self.quest_state_repository = quest_state_repository

    def get_all(
        self,
        db: Session,
        *,
        quest_id: int | None = None,
        field_path: str | None = None,
        field_path__startswith: str | None = None,
        automated: bool | None = None,
        recorded_at__gte: datetime | None = None,
        recorded_at__lte: datetime | None = None,
        latest: bool = False,
        ordering: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
        state_id: int | None = None,
    ) -> list[Record]:
        if state_id is not None:
            if recorded_at__gte is not None or recorded_at__lte is not None:
                raise BadRequestError(
                    "Cannot use 'state_id' together with 'recorded_at__gte' or 'recorded_at__lte'. "
                    "Please choose one filtering method."
                )

            state = self.quest_state_repository.get_by_id(db, state_id)
            if state is None:
                raise NotFoundError("state", state_id)

            recorded_at__gte = state.start_date
            recorded_at__lte = state.end_date

        if latest and (ordering or limit or offset):
            raise ValueError(
                "'latest' cannot be combined with 'ordering', 'limit', or 'offset'"
            )

        if recorded_at__gte and recorded_at__lte and recorded_at__gte > recorded_at__lte:
            raise ValueError("'recorded_at__gte' must be <= 'recorded_at__lte'")

        if (
            field_path
            and field_path__startswith
            and not field_path.startswith(field_path__startswith)
        ):
            raise ValueError("'field_path' must start with 'field_path__startswith'")

        records = self.record_repository.get_all(
            db,
            quest_id=quest_id,
            field_path=field_path,
            field_path__startswith=field_path__startswith,
            automated=automated,
            recorded_at__gte=recorded_at__gte,
            recorded_at__lte=recorded_at__lte,
            latest=latest,
            ordering=ordering,
            limit=limit,
            offset=offset,
        )

        logger.debug(f"Fetched {len(records)} records")
        return records

    def get(self, db: Session, record_id: int) -> Record:
        record = self.record_repository.get(db, record_id)
        if record is None:
            raise NotFoundError("Record", record_id)
        logger.debug(f"Record {record_id} found")
        return record

    def create(self, db: Session, request: RecordCreate, user_id: int) -> Record:
        quest = self.quest_repository.get_by_id(db, request.quest_id)

        if not quest:
            raise NotFoundError("Quest", request.quest_id)

        if quest.owner_id != user_id:
            raise ForbiddenError()

        quest_state = self.quest_state_repository.get_by_date(
            db, request.quest_id, request.recorded_at
        )

        if not quest_state:
            raise NotFoundError("Quest state")

        if not quest_state.fields:
            raise NoFieldsQuestStateError(quest_state.id)

        form_processor = FormProcessor()
        form_processor.load_fields(quest_state.fields)
        errors = form_processor.validate_value(request.field_path, value=request.value)

        if errors:
            raise RecordValidationFailed(request.value, request.field_path, errors)

        # Make an event that we want to create a new record
        # event = NewRecordEvent(field_id=request.field_path, new_value=request.value, recorded_at=request.recorded_at)
        # QuestService.handle_event(db, request.quest_id, event)

        # TODO: If allowed create new record
        record = self.record_repository.add(db, request)
        db.commit()
        logger.debug(f"Created record with id={record.id}")
        return record

    def remove(self, db: Session, record_id: int) -> None:
        record = self.record_repository.get(db, record_id)
        if record is None:
            raise NotFoundError("record", record_id)

        self.record_repository.remove(db, record_id)
        logger.debug(f"Removed record {record_id}")
