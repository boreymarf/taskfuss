from contextlib import contextmanager
import logging
from pathlib import Path

import sqlalchemy
from sqlalchemy.engine import Engine
from sqlalchemy.orm import Session, sessionmaker

from src.config import get_config
from src.db.base import Base

logger = logging.getLogger(__name__)


class Database:
    def __init__(self, db_url: str, echo: bool = False) -> None:
        self._engine = sqlalchemy.create_engine(db_url, echo=echo)
        self._session_factory = sessionmaker(
            bind=self._engine,
            autocommit=False,
            autoflush=False,
            expire_on_commit=False,
        )

    @classmethod
    def from_config(cls, echo: bool = False) -> "Database":
        db_path = get_config().database.db_path

        if db_path == ":memory:":
            db_url = "sqlite:///:memory:"
            logger.info("The engine was created in :memory:.")
        else:
            path = Path(db_path).resolve()
            path.parent.mkdir(parents=True, exist_ok=True)
            db_url = f"sqlite:///{path}"
            logger.info(f"The engine was created for '{path}' file")

        database = cls(db_url, echo=echo)
        database.create_database()
        return database

    @property
    def engine(self) -> Engine:
        return self._engine

    def create_database(self) -> None:
        Base.metadata.create_all(self._engine)

    @contextmanager
    def session(self):
        session: Session = self._session_factory()
        try:
            yield session
        except Exception:
            session.rollback()
            raise
        finally:
            session.close()

    def dispose(self) -> None:
        self._engine.dispose()
