from datetime import datetime
import logging
from sqlalchemy import or_
from sqlalchemy.orm import Session

from src.db.quest import QuestDB
from src.db.quest_state import QuestStateDB
from src.domain.quest_state import QuestState, QuestStateCreate

logger = logging.getLogger(__name__)

class QuestStateRepository:
    def add(self, db: Session, state_data: QuestStateCreate) -> QuestState:
        state_db = QuestStateDB(**state_data.model_dump())
        db.add(state_db)
        db.flush()
        db.refresh(state_db)
        logger.debug(f"Created new quest state for quest '{state_data.quest_id}'")
        return QuestState.model_validate(state_db)

    def get_by_id(self, db: Session, state_id: int) -> QuestState | None:
        state_db = db.query(QuestStateDB).filter(QuestStateDB.id == state_id).first()
        if state_db is None:
            return None
        return QuestState.model_validate(state_db)

    def get_latest(self, db: Session, quest_id: int) -> QuestState | None:
        state_db = (
            db.query(QuestStateDB)
            .filter(QuestStateDB.quest_id == quest_id)
            .order_by(QuestStateDB.start_date.desc())
            .first()
        )
        if state_db is None:
            return None
        return QuestState.model_validate(state_db)

    def has_overlap(
        self,
        db: Session,
        quest_id: int,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
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

    def get_all(
        self,
        db: Session,
        owner_id: int | None = None,
        quest_id: int | None = None,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
    ) -> list[QuestState]:
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

        states_db = query.all()
        logger.debug(
            f"Found {len(states_db)} states for owner_id={owner_id}, "
            f"quest_id={quest_id}, start_date={start_date}, end_date={end_date}"
        )
        return [QuestState.model_validate(state) for state in states_db]

    def get_by_date(
        self,
        db: Session,
        quest_id: int,
        target_date: datetime,
    ) -> QuestState | None:
        """
        Get the state active for the given quest at the exact target_date.
        Treats NULL start_date as 'beginning of time' and NULL end_date as 'end of time'.
        Returns the most recently started state if multiple overlap.
        """
        states_db = (
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

        if not states_db:
            logger.debug(f"No state found for quest {quest_id} at {target_date}")
            return None

        if len(states_db) > 1:
            logger.error(
                f"Multiple states ({len(states_db)}) found for quest {quest_id} at {target_date}. "
                "Returning the most recent (by start_date)."
            )
            return QuestState.model_validate(states_db[0])

        return QuestState.model_validate(states_db[0])
