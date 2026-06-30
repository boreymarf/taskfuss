from contextlib import contextmanager
import logging
from pathlib import Path
import sqlalchemy
from sqlalchemy.engine import Engine
from sqlalchemy.orm import Session
from src.config import get_config
from src.db.base import Base

logger = logging.getLogger(__name__)

_engine: Engine | None = None


def create_engine() -> Engine:
    global _engine

    if _engine is not None:
        logger.warning("Engine already created, returning existing.")
        return _engine

    try:
        if get_config().database.db_path == ":memory:":
            _engine = sqlalchemy.create_engine("sqlite:///:memory:")
            logger.info("The engine was created in :memory:.")
        else:
            db_path = Path(get_config().database.db_path).resolve()
            db_path.parent.mkdir(parents=True, exist_ok=True)

            logger.info(f"The engine was created for '{db_path}' file")
            _engine = sqlalchemy.create_engine(f"sqlite:///{db_path}")

        Base.metadata.create_all(_engine)
        return _engine

    except Exception as e:
        logger.critical(f"Failed to create engine: {e}")
        raise


def get_engine() -> Engine:
    global _engine
    if _engine is None:
        _engine = create_engine()
    return _engine


@contextmanager
def get_session():
    engine = get_engine()
    session = Session(engine)
    try:
        yield session
    finally:
        session.close()


def reset_engine():
    global _engine
    if _engine is not None:
        _engine.dispose()
        _engine = None
