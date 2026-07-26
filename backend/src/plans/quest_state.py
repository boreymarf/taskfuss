from datetime import datetime, timezone
from typing import Any
from uuid import UUID

from pydantic import BaseModel, ConfigDict, Field

from src.plans.fields import FormFields


class QuestState(BaseModel):
    id: UUID | None = None
    quest_id: UUID | None = None
    title: str | None = None
    start_date: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    end_date: datetime | None = None
    fields: FormFields | None = None
    data: dict[str, Any] | None = None

    model_config = ConfigDict(from_attributes=True)
