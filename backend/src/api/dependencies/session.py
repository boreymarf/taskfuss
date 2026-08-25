from fastapi import Request
from sqlalchemy.orm import Session


def get_session(request: Request) -> Session:
    """Session generator. Only used with FastAPI dependencies!"""
    database = request.app.state.database
    with database.session() as session:
        yield session
