from datetime import datetime
from typing import Any
from uuid import UUID

from pydantic import BaseModel, ConfigDict

from src.plans.fields import FormFields


class QuestStateCreate(BaseModel):
    quest_id: UUID
    title: str | None = None
    start_date: datetime | None = None
    end_date: datetime | None = None
    fields: FormFields | None = None
    data: dict[str, Any] | None = None

    model_config = ConfigDict(from_attributes=True)


class QuestState(BaseModel):
    id: UUID | None = None
    quest_id: UUID | None = None
    title: str | None = None
    start_date: datetime | None = None
    end_date: datetime | None = None
    fields: FormFields | None = None
    data: dict[str, Any] | None = None

    model_config = ConfigDict(from_attributes=True)
