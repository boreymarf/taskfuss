from abc import ABC, abstractmethod
from typing import Any

from src.quests.fields import Field
from src.quests.meta import PlanMeta


class BaseQuestPlan(ABC):
    """Abstract contract for a quest type module."""

    @abstractmethod
    def get_plan_info(self) -> PlanMeta:
        """Return metadata about this quest type (name, description, icon, etc.)."""
        ...

    @abstractmethod
    def get_setup_form(self) -> dict[str, Field]:
        """Return field definitions for the creation form."""
        ...

    def validate_setup(self, data: dict[str, Any]) -> list[str]:
        """Validate creation form data. Returns list of errors (empty if valid)."""
        return []

    @abstractmethod
    def get_quest_form(self, current_state: dict[str, Any]) -> dict[str, Field]:
        """Return field definitions for the interaction form, based on current state."""
        ...

    def validate_record(
        self, field_id: str, new_value: Any, current_state: dict[str, Any]
    ) -> list[str]:
        """Validate a single field change before persisting. Returns list of errors."""
        ...
        return []

    def on_record(
        self, field_id: str, new_value: Any, current_state: dict[str, Any]
    ) -> list[dict[str, Any]]:
        """Called after a record is persisted. Returns list of side-effect actions."""
        ...

    def on_event(self, event_name: str, current_state: dict[str, Any]) -> list[dict[str, Any]]:
        """Handle system events (e.g. 'daily_reset'). Returns list of side-effect actions."""
        ...
