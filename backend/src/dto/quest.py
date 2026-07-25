from typing import Any

from pydantic import BaseModel, ConfigDict


class QuestCreate(BaseModel):

    plan_id: str
    setup_form_data: dict[str, Any]

    model_config = ConfigDict(from_attributes=True)
