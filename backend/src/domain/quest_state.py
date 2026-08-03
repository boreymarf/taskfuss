from datetime import datetime
from typing import Any

from pydantic import BaseModel, ConfigDict

from src.domain.fields import FormFields


class QuestStateCreate(BaseModel):
    quest_id: int
    title: str | None = None
    start_date: datetime | None = None
    end_date: datetime | None = None
    fields: FormFields | None = None
    data: dict[str, Any] | None = None

    model_config = ConfigDict(from_attributes=True)


class QuestState(BaseModel):
    id: int | None = None
    quest_id: int | None = None
    title: str | None = None
    start_date: datetime | None = None
    end_date: datetime | None = None
    fields: FormFields | None = None
    data: dict[str, Any] | None = None

    model_config = ConfigDict(from_attributes=True)


class QuestStateResponse(BaseModel):
    id: int
    quest_id: int
    title: str | None
    start_date: datetime
    end_date: datetime | None
    fields: FormFields | None
    data: dict[str, Any] | None

    model_config = ConfigDict(from_attributes=True)
