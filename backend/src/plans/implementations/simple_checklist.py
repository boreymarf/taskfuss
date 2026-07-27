from typing import Any, override

from src.domain import (
    FormFields,
    InitEvent,
    ListField,
    PlanMetadata,
    QuestAction,
    StrField,
)
from src.plans.base import BasePlan


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
