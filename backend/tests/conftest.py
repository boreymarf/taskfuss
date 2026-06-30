import os
from pathlib import Path

import pytest
from sqlalchemy import create_engine
from sqlalchemy.orm import Session

from src.db import Base


@pytest.fixture(scope="function")
def project_root():
    return Path(os.getcwd())


@pytest.fixture(scope="function")
def db_session():
    engine = create_engine("sqlite://", echo=False)
    Base.metadata.create_all(engine)

    with Session(engine) as session:
        yield session

    engine.dispose()


@pytest.fixture(scope="class")
def db_session_class():
    engine = create_engine("sqlite://", echo=False)
    Base.metadata.create_all(engine)

    with Session(engine) as session:
        yield session

    engine.dispose()
