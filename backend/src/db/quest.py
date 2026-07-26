from typing import Any
from uuid import UUID, uuid4

from sqlalchemy import JSON, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from src.db import Base


class QuestDB(Base):
    __tablename__ = "quest"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    owner_id: Mapped[int] = mapped_column(ForeignKey("user.id", ondelete="CASCADE"))
    plan_id: Mapped[str] = mapped_column(
        ForeignKey("plan.id", ondelete="CASCADE")
    )
