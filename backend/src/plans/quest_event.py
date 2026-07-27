from typing import Any, Literal

from pydantic import BaseModel, ConfigDict


class QuestEvent(BaseModel):
    discriminator: str

    model_config = ConfigDict(from_attributes=True)


class InitEvent(QuestEvent):
    discriminator: Literal["init"] = "init"


class NewRecordEvent(QuestEvent):
    discriminator: Literal["new_record"] = "new_record"
    field_id: str
    new_value: Any


class SetupUpdateEvent(QuestEvent):
    discriminator: Literal["setup_update"] = "setup_update"
