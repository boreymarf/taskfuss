from __future__ import annotations
from datetime import datetime
from typing import Any, Literal
from pydantic import BaseModel, ConfigDict

from src.plans.quest_data import QuestData


class Action(BaseModel):
    discriminator: str
    model_config = ConfigDict(from_attributes=True)


class UpdateStateAction(Action):
    discriminator: Literal["update_state"] = "update_state"
    new_data: QuestData


class CreateNewStateAction(Action):
    discriminator: Literal["create_new_state"] = "create_new_state"
    data: QuestData
    creation_date: datetime


class ScheduleEventAction(Action):
    discriminator: Literal["schedule_event"] = "schedule_event"
    event_name: str
    trigger: Literal["once", "interval"] = "once"
    interval_seconds: int | None = None


class UnscheduleEventAction(Action):
    discriminator: Literal["unschedule_event"] = "unschedule_event"
    event_name: str


class AddRecordAction(Action):
    discriminator: Literal["add_record"] = "add_record"
    field_path: str
    value: Any
