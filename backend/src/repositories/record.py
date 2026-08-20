from datetime import datetime

from sqlalchemy import ColumnElement, func, desc, asc
from sqlalchemy.orm import Session

from src.db.record import RecordDB
from src.domain.record import RecordQueryParams



class RecordRepository:
    @staticmethod
    def add(db: Session, record: RecordDB) -> RecordDB:
        db.add(record)
        db.flush()
        db.refresh(record)
        return record

    @staticmethod
    def remove(db: Session, record: RecordDB) -> None:
        db.delete(record)
        db.flush()

    @staticmethod
    def get(db: Session, record_id: int) -> RecordDB | None:
        return db.get(RecordDB, record_id)

    @staticmethod
    def get_all(
        db: Session,
        *,
        quest_id: int | None = None,
        field_path: str | None = None,
        field_path__startswith: str | None = None,
        automated: bool | None = None,
        created_at__gte: datetime | None = None,
        created_at__lte: datetime | None = None,
        latest: bool = False,
        ordering: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
    ) -> list[RecordDB]:
        filters: list[ColumnElement[bool]] = []

        if quest_id is not None:
            filters.append(RecordDB.quest_id == quest_id)
        if field_path is not None:
            filters.append(RecordDB.field_path == field_path)
        if field_path__startswith is not None:
            filters.append(RecordDB.field_path.startswith(field_path__startswith))
        if automated is not None:
            filters.append(RecordDB.automated == automated)
        if created_at__gte is not None:
            filters.append(RecordDB.created_at >= created_at__gte)
        if created_at__lte is not None:
            filters.append(RecordDB.created_at <= created_at__lte)

        query = db.query(RecordDB)
        for condition in filters:
            query = query.filter(condition)

        if latest:
            subquery = db.query(
                RecordDB.field_path,
                func.max(RecordDB.created_at).label("max_created"),
            )
            for condition in filters:
                subquery = subquery.filter(condition)
            subquery = subquery.group_by(RecordDB.field_path).subquery()

            query = db.query(RecordDB).join(
                subquery,
                (RecordDB.field_path == subquery.c.field_path)
                & (RecordDB.created_at == subquery.c.max_created),
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

        return query.all()
