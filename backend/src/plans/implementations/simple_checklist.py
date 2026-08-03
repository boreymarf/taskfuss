from pprint import pprint
from typing import override

from src.domain import (
    FormFields,
    InitEvent,
    ListField,
    PlanMetadata,
    QuestAction,
    QuestStateCreate,
    StrField,
)
from src.domain.fields import CheckboxField, TupleField
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

        quest_form: FormFields = {}

        tasks = settings.form_data.get("task_list")
        checkbox_fields: list[CheckboxField] = []
        if tasks:
            for task in tasks:
                checkbox_fields.append(CheckboxField(label=task))

        quest_form["task_list"] = TupleField(fields=tuple(checkbox_fields))

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
