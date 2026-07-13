from enum import Enum
from uuid import UUID

from pydantic import BaseModel, ConfigDict


class QuestState(Enum):
    ACTIVE = "active"
    PAUSED = "paused"
    COMPLETED = "completed"
    DROPPED = "dropped"


class Quest(BaseModel):
    """All info about current instance of the quest"""

    id: UUID
    title: str | None = None
    completion: float = 0
    state: QuestState = QuestState.ACTIVE
    plan_id: str
    setup_form_data: str

    model_config = ConfigDict(from_attributes=True)
