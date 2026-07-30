from typing import Any
from uuid import UUID
from pydantic import BaseModel, ConfigDict


class QuestCreateRequest(BaseModel):
    """DTO"""

    plan_id: str
    settings: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)


class QuestCreateResponse(BaseModel):
    """DTO"""

    id: UUID
    owner_id: int
    plan_id: str

    model_config = ConfigDict(from_attributes=True)


class QuestCreate(BaseModel):
    """Internal struct for repo"""

    owner_id: int
    plan_id: str
    settings: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)


class Quest(BaseModel):
    """All info about current instance of the quest"""

    id: UUID
    owner_id: int
    plan_id: str

    model_config = ConfigDict(from_attributes=True)
