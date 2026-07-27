from typing import Any, override

from src.plans.base import BasePlan
from src.plans.fields import CheckboxField, FormFields, ListField, StrField
from src.plans.metadata import PlanMetadata
from src.plans.quest_actions import QuestAction
from src.plans.quest_event import InitEvent
from src.plans.quest_state import QuestStateCreate


class SimpleChecklist(BasePlan):

    @override
    def get_metadata(self) -> PlanMetadata:
        return PlanMetadata(
            id="simple_checklist",
            name="Simple Checklist",
            description="This is a simple checklist",
        )

    @override
    def get_setup_fields(self) -> FormFields:
        return {
            "name": StrField(label="Quest name (non optional)", required=True),
            "task_list": ListField(item_field=StrField()),
        }

    @override
    def on_init(self, event: InitEvent) -> list[QuestAction]:
        return []

