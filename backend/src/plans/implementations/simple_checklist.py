from typing import Any, override

from src.plans.base import BasePlan
from src.plans.fields import FormFields, ListField, StrField
from src.plans.metadata import PlanMetadata
from src.plans.quest_data import QuestData
from src.plans.setup_data import SetupData


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
    def validate_setup_data(self, setup_info: SetupData) -> list[str]:
        return []

    @override
    def compute_initial_quest_data(self, setup_data: SetupData) -> QuestData:
        return QuestData()

    @override
    def compute_current_quest_data(
        self, current_data: QuestData, setup_data: SetupData
    ) -> QuestData:

        return QuestData()
