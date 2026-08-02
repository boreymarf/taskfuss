from uuid import UUID
from sqlalchemy import select
from sqlalchemy.orm import Session

from src.db.record import RecordDB



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
    def get(db: Session, record_id: UUID) -> RecordDB | None:
        return db.get(RecordDB, record_id)

    @staticmethod
    def get_all(db: Session) -> list[RecordDB]:
        return db.query(RecordDB).all()

    @staticmethod
    def get_all_by_quest_id(db: Session, quest_id: UUID) -> list[RecordDB]:
        return db.query(RecordDB).filter(RecordDB.quest_id == quest_id).all()
