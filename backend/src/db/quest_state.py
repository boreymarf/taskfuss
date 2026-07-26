from datetime import datetime, timezone
from typing import Any
from uuid import UUID, uuid4

from sqlalchemy import JSON, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column

from src.db import Base
from src.plans.fields import FormFields


class QuestStateDB(Base):
    __tablename__ = "quest_state"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    quest_id: Mapped[UUID] = mapped_column(ForeignKey("quest.id", ondelete="CASCADE"))
    title: Mapped[str | None] = mapped_column(nullable=True)
    start_date: Mapped[datetime] = mapped_column(default=datetime.now(timezone.utc))
    end_date: Mapped[datetime | None] = mapped_column(nullable=True)
    fields: Mapped[FormFields | None] = mapped_column(JSON, nullable=True)
    data: Mapped[dict[str, Any] | None] = mapped_column(JSON, nullable=True)
