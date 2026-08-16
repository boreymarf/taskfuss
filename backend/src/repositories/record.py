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
    def get_all(db: Session, params: RecordQueryParams) -> list[RecordDB]:
        filters: list[ColumnElement[bool]] = []

        if params.quest_id is not None:
            filters.append(RecordDB.quest_id == params.quest_id)
        if params.field_path is not None:
            filters.append(RecordDB.field_path == params.field_path)
        if params.field_path__startswith is not None:
            filters.append(RecordDB.field_path.startswith(params.field_path__startswith))
        if params.automated is not None:
            filters.append(RecordDB.automated == params.automated)
        if params.created_at__gte is not None:
            filters.append(RecordDB.created_at >= params.created_at__gte)
        if params.created_at__lte is not None:
            filters.append(RecordDB.created_at <= params.created_at__lte)

        query = db.query(RecordDB)
        for condition in filters:
            query = query.filter(condition)

        if params.latest:
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

        if params.ordering:
            field_name = params.ordering.lstrip("-")
            if not hasattr(RecordDB, field_name):
                raise ValueError(f"Invalid ordering field: {field_name}")
            order_func = desc if params.ordering.startswith("-") else asc
            query = query.order_by(order_func(getattr(RecordDB, field_name)))

        if params.offset is not None:
            query = query.offset(params.offset)
        if params.limit is not None:
            query = query.limit(params.limit)

        return query.all()
