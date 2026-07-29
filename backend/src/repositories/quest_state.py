from uuid import UUID
from sqlalchemy.orm import Session

from src.db.quest_state import QuestStateDB
from src.domain import QuestStateCreate


class QuestStateRepository:
    @staticmethod
    def add(db: Session, state_data: QuestStateCreate) -> QuestStateDB:
        state_db = QuestStateDB(**state_data.model_dump())
        db.add(state_db)
        db.flush()
        db.refresh(state_db)
        return state_db

    @staticmethod
    def get_latest(db: Session, quest_id: UUID) -> QuestStateDB | None:
        return (
            db.query(QuestStateDB)
            .filter(QuestStateDB.quest_id == quest_id)
            .order_by(QuestStateDB.start_date.desc())
            .first()
        )
