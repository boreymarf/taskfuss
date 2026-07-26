from uuid import UUID
from datetime import datetime
from typing import Any
from pydantic import BaseModel, ConfigDict

class Record(BaseModel):
    id: UUID
    quest_id: UUID
    state_id: UUID | None = None
    field_path: str
    value: Any
    created_at: datetime
    user_id: int | None = None

    model_config = ConfigDict(from_attributes=True)
