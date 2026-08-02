from uuid import UUID
from datetime import datetime
from typing import Any
from pydantic import BaseModel, ConfigDict


class Record(BaseModel):
    """Domain level model"""
    id: UUID
    automated: bool = False
    quest_id: UUID
    field_path: str
    value: Any
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)


class RecordCreateRequest(BaseModel):
    """DTO"""

    field_path: str
    quest_id: UUID
    value: Any

    model_config = ConfigDict(from_attributes=True)


class RecordResponse(BaseModel):
    """DTO"""

    id: UUID
    created_at: datetime

    field_path: str
    quest_id: UUID
    automated: bool
    value: Any

    model_config = ConfigDict(from_attributes=True)


class RecordCreate(BaseModel):
    """Used by repo"""

    field_path: str
    quest_id: UUID
    value: Any
    automated: bool = False

    model_config = ConfigDict(from_attributes=True)
