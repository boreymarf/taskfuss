from datetime import datetime
import logging
from uuid import UUID
from sqlalchemy import or_
from sqlalchemy.orm import Session

from src.db.quest import QuestDB
from src.db.quest_state import QuestStateDB
from src.domain import QuestStateCreate

logger = logging.getLogger(__name__)

class QuestStateRepository:
    @staticmethod
    def add(db: Session, state_data: QuestStateCreate) -> QuestStateDB:
        state_db = QuestStateDB(**state_data.model_dump())
        db.add(state_db)
        db.flush()
        db.refresh(state_db)
        logger.debug(f"Created new quest state for quest '{state_data.quest_id}'")
        return state_db

    @staticmethod
    def get_latest(db: Session, quest_id: UUID) -> QuestStateDB | None:
        return (
            db.query(QuestStateDB)
            .filter(QuestStateDB.quest_id == quest_id)
            .order_by(QuestStateDB.start_date.desc())
            .first()
        )

    @staticmethod
    def has_overlap(
        db: Session,
        quest_id: UUID,
        start_date: datetime | None = None,
        end_date: datetime | None = None
    ) -> bool:
        """
        Return True if any existing state overlaps with [start_date, end_date).
        None boundary means unbounded interval.
        """
        conditions = [QuestStateDB.quest_id == quest_id]

        if end_date is not None:
            conditions.append(QuestStateDB.start_date < end_date)

        if start_date is not None:
            conditions.append(
                or_(
                    QuestStateDB.end_date.is_(None),
                    QuestStateDB.end_date > start_date,
                )
            )

        return db.query(
            db.query(QuestStateDB).filter(*conditions).exists()
        ).scalar()

    @staticmethod
    def get_all_by_date(db: Session, target_date: datetime, owner_id: int) -> list[QuestStateDB]:
        states = (
            db.query(QuestStateDB)
            .join(QuestDB, QuestDB.id == QuestStateDB.quest_id)
            .filter(
                QuestDB.owner_id == owner_id,
                QuestStateDB.start_date <= target_date,
                or_(
                    QuestStateDB.end_date.is_(None),
                    QuestStateDB.end_date > target_date
                )
            )
            .all()
        )
        logger.debug(f"Found {len(states)} states for owner_id={owner_id} at {target_date}")
        return states
