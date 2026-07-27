from typing import Any

from sqlalchemy import JSON, ForeignKey
from sqlalchemy.orm import Mapped, mapped_column
from uuid import UUID, uuid4
from datetime import datetime, timezone

from src.db import Base


class QuestSettingsDB(Base):
    __tablename__ = "quest_settings"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)
    quest_id: Mapped[UUID] = mapped_column(ForeignKey("quest.id", ondelete="CASCADE"))
    form_data: Mapped[dict[str, Any]] = mapped_column(JSON, nullable=False)

    created_at: Mapped[datetime] = mapped_column(default=datetime.now(timezone.utc))
