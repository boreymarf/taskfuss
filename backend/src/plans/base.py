from abc import ABC, abstractmethod
from typing import Any

from src.plans.fields import FormFields
from src.plans.metadata import PlanMetadata
from src.plans.quest_actions import QuestAction
from src.plans.quest_event import InitEvent, NewRecordEvent, QuestEvent
from src.plans.quest_settings import QuestSettingsCreate
from src.plans.quest_state import QuestStateCreate


class BasePlan(ABC):
    """Abstract contract for a quest type module."""

    @abstractmethod
    def get_metadata(self) -> PlanMetadata:
        """Return metadata about this quest type (name, description, icon, etc.)."""
        ...

    @abstractmethod
    def get_setup_fields(self) -> FormFields:
        """Return field definitions for the creation form."""
        ...

    def validate_settings_form(self, _form: dict[str, Any]) -> list[str]:
        """Validate creation form data. Returns list of errors (empty if valid)."""
        return []

    def handle_event(self, event: QuestEvent) -> list[QuestAction]:
        if isinstance(event, InitEvent):
            return self.on_init(event)
        if isinstance(event, NewRecordEvent):
            return self.on_new_record_event(event)
        return self.on_custom_event(event)

    def on_init(self, _event: InitEvent) -> list[QuestAction]:
        """When quest is created. Must declare at least one state."""
        return []

    def on_new_record_event(self, _event: NewRecordEvent) -> list[QuestAction]:
        """When non-automatic (user's) record is added."""
        return []

    def on_custom_event(self, _event: QuestEvent) -> list[QuestAction]:
        """Anything else"""
        return []
