from datetime import datetime
import logging
from sqlalchemy.orm import Session

from src.classes.form_processor import FormProcessor
from src.db.record import RecordDB
from src.domain.record import Record, RecordCreate
from src.exceptions.generic import BadRequestError, ForbiddenError, NotFoundError
from src.exceptions.quest_state import NoFieldsQuestStateError
from src.exceptions.record import RecordValidationFailed
from src.repositories.quest import QuestRepository
from src.repositories.quest_state import QuestStateRepository
from src.repositories.record import RecordRepository

logger = logging.getLogger(__name__)


class RecordService:
    @staticmethod
    def get_all(
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

            state = QuestStateRepository.get_by_id(db, state_id)
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

        records_db = RecordRepository.get_all(
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

        result = [Record.model_validate(r) for r in records_db]
        logger.debug(f"Fetched {len(result)} records")
        return result

    @staticmethod
    def get(db: Session, record_id: int) -> Record | None:
        record_db = RecordRepository.get(db, record_id)
        if record_db is None:
            raise NotFoundError("Record", record_id)
        result = Record.model_validate(record_db)
        logger.debug(f"Record {record_id} found")
        return result

    @staticmethod
    def create(db: Session, request: RecordCreate, user_id: int) -> Record:
        quest_db = QuestRepository.get_by_id(db, request.quest_id)

        if not quest_db:
            raise NotFoundError("Quest", request.quest_id)

        if quest_db.owner_id != user_id:
            raise ForbiddenError()

        quest_state_db = QuestStateRepository.get_by_date(
            db, request.quest_id, request.recorded_at
        )

        if not quest_state_db:
            raise NotFoundError("Quest state")

        if not quest_state_db.fields:
            raise NoFieldsQuestStateError(quest_state_db.id)

        form_processor = FormProcessor()
        form_processor.load_fields(quest_state_db.fields)
        errors = form_processor.validate_value(request.field_path, value=request.value)

        if errors:
            raise RecordValidationFailed(request.value, request.field_path, errors)

        # TODO: I don't think it should be here
        record_db = RecordDB(**request.model_dump())
        created = RecordRepository.add(db, record_db)
        db.commit()
        result = Record.model_validate(created)
        logger.debug(f"Created record with id={result.id}")
        return result

    @staticmethod
    def remove(db: Session, record_id: int) -> None:
        record_db = RecordRepository.get(db, record_id)
        if record_db is None:
            raise NotFoundError("record", record_id)

        RecordRepository.remove(db, record_db)
        logger.debug(f"Removed record {record_id}")
