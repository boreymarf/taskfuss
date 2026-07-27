from __future__ import annotations
from datetime import datetime
from typing import Literal
from pydantic import BaseModel, ConfigDict

from src.plans.quest_state import QuestStateCreate



class QuestActionBase(BaseModel):
    discriminator: str

    model_config = ConfigDict(from_attributes=True)


class UpdateStateAction(QuestActionBase):
    discriminator: Literal["update_state"] = "update_state"
    data: QuestStateCreate


class CreateNewStateAction(QuestActionBase):
    discriminator: Literal["create_new_state"] = "create_new_state"
    data: QuestStateCreate
    creation_date: datetime


class ScheduleEventAction(QuestActionBase):
    discriminator: Literal["schedule_event"] = "schedule_event"
    event_name: str
    trigger: Literal["once", "interval"] = "once"
    interval_seconds: int | None = None


# class UnscheduleEventAction(QuestActionBase):
#     discriminator: Literal["unschedule_event"] = "unschedule_event"
#     event_name: str
#
#
# class AddRecordAction(QuestActionBase):
#     discriminator: Literal["add_record"] = "add_record"
#     field_path: str
#     value: Any

QuestAction = UpdateStateAction | CreateNewStateAction | ScheduleEventAction
