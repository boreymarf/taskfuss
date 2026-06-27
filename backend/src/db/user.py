from sqlalchemy import String
from sqlalchemy.orm import Mapped, mapped_column
from src.db.base import Base


class UserDB(Base):
    __tablename__ = "user"

    id: Mapped[int] = mapped_column(primary_key=True)
    login: Mapped[str] = mapped_column(String(30), unique=True)
    password_hash: Mapped[str] = mapped_column(String(128))
