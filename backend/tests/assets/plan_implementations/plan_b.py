from typing import Any, override

from src.quests.base import BaseQuestPlan
from src.quests.fields import Field, ListField, StrField
from src.quests.meta import PlanMeta


class PlanA(BaseQuestPlan):
    @override
    def get_plan_info(self) -> PlanMeta:
        return PlanMeta(
            id="plan_b",
            name="name_b",
            description="desc_b",
        )

    @override
    def get_setup_form(self) -> dict[str, Field]:
        return {
            "name": StrField(label="Quest name (optional)"),
            "task_list": ListField(item_field=StrField()),
        }

    @override
    def get_quest_form(self, current_state: dict[str, Any]) -> dict[str, Field]:
        return {}
