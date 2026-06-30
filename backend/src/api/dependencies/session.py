from sqlalchemy.orm import Session

from src.database import get_engine


def get_session():
    """Session generator. Only used with FastAPI dependencies!"""
    engine = get_engine()
    with Session(engine) as session:
        yield session
