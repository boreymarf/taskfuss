from datetime import datetime
from typing import Any
from fastapi import Query
from pydantic import BaseModel, ConfigDict, Field, model_validator


class Record(BaseModel):
    """Domain level model"""

    id: int
    automated: bool = False
    quest_id: int
    field_path: str
    value: Any
    created_at: datetime

    model_config = ConfigDict(from_attributes=True)


class RecordCreateRequest(BaseModel):
    """DTO"""

    field_path: str
    quest_id: int
    value: Any

    model_config = ConfigDict(from_attributes=True)


class RecordResponse(BaseModel):
    """DTO"""

    id: int
    created_at: datetime

    field_path: str
    quest_id: int
    automated: bool
    value: Any

    model_config = ConfigDict(from_attributes=True)


class RecordCreate(BaseModel):
    """Used by repo"""

    field_path: str
    quest_id: int
    value: Any
    automated: bool = False

    model_config = ConfigDict(from_attributes=True)


class RecordQueryParams(BaseModel):
    quest_id: int | None = None
    field_path: str | None = None
    field_path__startswith: str | None = None
    automated: bool | None = None
    created_at__gte: datetime | None = None
    created_at__lte: datetime | None = None
    latest: bool = False
    ordering: str | None = None
    limit: int | None = Field(None, ge=1, description="Maximum number of records")
    offset: int | None = Field(None, ge=0, description="Number of records to skip")

    @model_validator(mode="after")
    def check_conflicts(self):
        if self.latest and (self.ordering or self.limit or self.offset):
            raise ValueError("'latest' cannot be combined with 'ordering', 'limit', or 'offset'")
        if self.created_at__gte and self.created_at__lte and self.created_at__gte > self.created_at__lte:
            raise ValueError("'created_at__gte' must be <= 'created_at__lte'")
        if self.field_path and self.field_path__startswith and not self.field_path.startswith(self.field_path__startswith):
            raise ValueError("'field_path' must start with 'field_path__startswith'")
        return self
