from enum import Enum
from uuid import UUID

from pydantic import BaseModel, ConfigDict

# class QuestStatus(Enum):
#     ACTIVE = "active"
#     PAUSED = "paused"
#     COMPLETED = "completed"
#     DROPPED = "dropped"
#


class QuestCreate(BaseModel):
    owner_id: int
    plan_id: str

    model_config = ConfigDict(from_attributes=True)


class Quest(BaseModel):
    """All info about current instance of the quest"""

    id: UUID
    owner_id: int
    plan_id: str

    model_config = ConfigDict(from_attributes=True)
