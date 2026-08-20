from datetime import datetime
import logging
from sqlalchemy import or_
from sqlalchemy.orm import Session

from src.db.quest import QuestDB
from src.db.quest_state import QuestStateDB
from src.domain.quest_state import QuestStateCreate

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
    def get_by_id(db: Session, state_id: int) -> QuestStateDB | None:
        return db.query(QuestStateDB).filter(QuestStateDB.id == state_id).first()

    @staticmethod
    def get_latest(db: Session, quest_id: int) -> QuestStateDB | None:
        return (
            db.query(QuestStateDB)
            .filter(QuestStateDB.quest_id == quest_id)
            .order_by(QuestStateDB.start_date.desc())
            .first()
        )

    @staticmethod
    def has_overlap(
        db: Session,
        quest_id: int,
        start_date: datetime | None = None,
        end_date: datetime | None = None
    ) -> bool:
        """
        Return True if any existing state overlaps with [start_date, end_date).
        None boundary means unbounded interval.
        """
        conditions = [QuestStateDB.quest_id == quest_id]

        if end_date is not None:
            conditions.append(
                or_(
                    QuestStateDB.start_date.is_(None),        
                    QuestStateDB.start_date < end_date
                )
            )

        if start_date is not None:
            conditions.append(
                or_(
                    QuestStateDB.end_date.is_(None),         
                    QuestStateDB.end_date > start_date
                )
            )

        return db.query(
            db.query(QuestStateDB).filter(*conditions).exists()
        ).scalar()

    @staticmethod
    def get_all(
        db: Session,
        owner_id: int | None = None,
        quest_id: int | None = None,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
    ) -> list[QuestStateDB]:
        """
        Get all quest states, optionally filtered by owner, quest and date range.
        Dates are compared strictly (>= for start, <= for end).
        NULL values are **excluded** from these filters (they don't satisfy >= or <=).
        """
        query = db.query(QuestStateDB)

        if owner_id is not None:
            query = query.join(QuestDB, QuestDB.id == QuestStateDB.quest_id)
            query = query.filter(QuestDB.owner_id == owner_id)

        if quest_id is not None:
            query = query.filter(QuestStateDB.quest_id == quest_id)

        if start_date is not None:
            query = query.filter(QuestStateDB.start_date >= start_date)

        if end_date is not None:
            query = query.filter(QuestStateDB.end_date <= end_date)

        states = query.all()
        logger.debug(
            f"Found {len(states)} states for owner_id={owner_id}, "
            f"quest_id={quest_id}, start_date={start_date}, end_date={end_date}"
        )
        return states

    @staticmethod
    def get_by_date(
        db: Session,
        quest_id: int,
        target_date: datetime,
    ) -> QuestStateDB | None:
        """
        Get the state active for the given quest at the exact target_date.
        Treats NULL start_date as 'beginning of time' and NULL end_date as 'end of time'.
        Returns the most recently started state if multiple overlap.
        """
        states = (
            db.query(QuestStateDB)
            .filter(
                QuestStateDB.quest_id == quest_id,
                or_(
                    QuestStateDB.start_date.is_(None),
                    QuestStateDB.start_date <= target_date
                ),
                or_(
                    QuestStateDB.end_date.is_(None),
                    QuestStateDB.end_date > target_date
                )
            )
            .order_by(QuestStateDB.start_date.desc().nulls_last())
            .limit(2)
            .all()
        )

        if not states:
            logger.debug(f"No state found for quest {quest_id} at {target_date}")
            return None

        if len(states) > 1:
            logger.error(
                f"Multiple states ({len(states)}) found for quest {quest_id} at {target_date}. "
                "Returning the most recent (by start_date)."
            )
            return states[0]

        return states[0]
