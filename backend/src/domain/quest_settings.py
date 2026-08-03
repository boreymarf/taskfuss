from typing import Any

from pydantic import BaseModel, ConfigDict


class QuestSettingsCreate(BaseModel):
    """All info gathered during setup"""

    form_data: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)


class QuestSettings(BaseModel):

    id: int
    quest_id: int
    form_data: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)
