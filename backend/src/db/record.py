from datetime import datetime, timezone
from typing import Any
from uuid import UUID, uuid4

from sqlalchemy import JSON, ForeignKey, String
from sqlalchemy.orm import Mapped, mapped_column

from src.db import Base


class RecordDB(Base):
    __tablename__ = "record"

    id: Mapped[UUID] = mapped_column(primary_key=True, default=uuid4)

    automated: Mapped[bool] = mapped_column(default=False)
    field_path: Mapped[str] = mapped_column(String, index=True)
    quest_id: Mapped[UUID] = mapped_column(ForeignKey("quest.id", ondelete="CASCADE"))
    value: Mapped[Any] = mapped_column(JSON)
    created_at: Mapped[datetime] = mapped_column(default=datetime.now(timezone.utc))

    user_id: Mapped[int] = mapped_column(ForeignKey("user.id"), nullable=True)
