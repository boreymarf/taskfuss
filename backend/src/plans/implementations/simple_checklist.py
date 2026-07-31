from pprint import pprint
from typing import override
import uuid

from src.domain import (
    FormFields,
    InitEvent,
    ListField,
    PlanMetadata,
    QuestAction,
    QuestStateCreate,
    StrField,
)
from src.domain.fields import CheckboxField
from src.domain.quest_actions import CreateNewStateAction
from src.plans.base import BasePlan
from src.plans.quest_context import QuestContext


class SimpleChecklist(BasePlan):

    @override
    def get_metadata(self) -> PlanMetadata:
        return PlanMetadata(
            id="simple_checklist",
            name="Simple Checklist",
            description="This is a simple checklist",
        )

    @override
    def get_settings_form(self) -> FormFields:
        return {
            "title": StrField(label="Quest name (non optional)", required=True),
            "task_list": ListField(item_field=StrField()),
        }

    @override
    def on_init(self, ctx: QuestContext, event: InitEvent) -> list[QuestAction]:
        settings = ctx.get_quest_settings()
        assert settings != None

        pprint(settings)

        quest_form: FormFields = {}

        if settings.form_data.get("task_list"):
            for i, item_name in enumerate(settings.form_data["task_list"]):
                quest_form[str(i)] = CheckboxField(label=item_name)

        actions: list[QuestAction] = [
            CreateNewStateAction(
                data=QuestStateCreate(
                    quest_id=ctx.get_quest_id(),
                    title=settings.form_data["title"],
                    fields=quest_form,
                )
            )
        ]
        return actions
