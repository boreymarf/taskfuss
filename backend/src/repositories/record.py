from datetime import datetime
from sqlalchemy import ColumnElement, func, desc, asc
from sqlalchemy.orm import Session

from src.db.record import RecordDB
from src.domain.record import Record, RecordCreate

import logging

logger = logging.getLogger(__name__)


class RecordRepository:
    def add(self, db: Session, record: RecordCreate) -> Record:
        record_db = RecordDB(**record.model_dump())
        db.add(record_db)
        db.flush()
        db.refresh(record_db)
        new_record = Record.model_validate(record_db)
        logger.debug(f"Added record with id={new_record.id}")
        return new_record

    def remove(self, db: Session, record_id: int) -> None:
        record_db = db.get(RecordDB, record_id)
        if record_db:
            db.delete(record_db)
            db.flush()
            logger.debug(f"Removed record {record_id}")

    def get(self, db: Session, record_id: int) -> Record | None:
        record_db = db.get(RecordDB, record_id)
        if record_db is None:
            return None
        return Record.model_validate(record_db)

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
    ) -> list[Record]:
        filters: list[ColumnElement[bool]] = []

        if quest_id is not None:
            filters.append(RecordDB.quest_id == quest_id)
        if field_path is not None:
            filters.append(RecordDB.field_path == field_path)
        if field_path__startswith is not None:
            filters.append(RecordDB.field_path.startswith(field_path__startswith))
        if automated is not None:
            filters.append(RecordDB.automated == automated)
        if recorded_at__gte is not None:
            filters.append(RecordDB.recorded_at >= recorded_at__gte)
        if recorded_at__lte is not None:
            filters.append(RecordDB.recorded_at <= recorded_at__lte)

        query = db.query(RecordDB)
        for condition in filters:
            query = query.filter(condition)

        if latest:
            subquery = db.query(
                RecordDB.field_path,
                func.max(RecordDB.recorded_at).label("max_created"),
            )
            for condition in filters:
                subquery = subquery.filter(condition)
            subquery = subquery.group_by(RecordDB.field_path).subquery()

            query = db.query(RecordDB).join(
                subquery,
                (RecordDB.field_path == subquery.c.field_path)
                & (RecordDB.recorded_at == subquery.c.max_created),
            )

        if ordering:
            field_name = ordering.lstrip("-")
            if not hasattr(RecordDB, field_name):
                raise ValueError(f"Invalid ordering field: {field_name}")
            order_func = desc if ordering.startswith("-") else asc
            query = query.order_by(order_func(getattr(RecordDB, field_name)))

        if offset is not None:
            query = query.offset(offset)
        if limit is not None:
            query = query.limit(limit)

        records_db = query.all()
        logger.debug(f"Fetched {len(records_db)} records from database")
        return [Record.model_validate(r) for r in records_db]
