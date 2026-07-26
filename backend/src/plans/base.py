from abc import ABC, abstractmethod
from typing import Any

from src.plans.fields import FormFields
from src.plans.metadata import PlanMetadata
from src.plans.quest_state import QuestState
from src.plans.setup_data import SetupData


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

    def validate_setup_data(self, _setup_info: SetupData) -> list[str]:
        """Validate creation form data. Returns list of errors (empty if valid)."""
        return []

    def compute_initial_quest_state(self, setup_data: SetupData) -> QuestState:
        """Return field definitions for the interaction form, based on current state."""
        ...

    def compute_current_quest_data(
        self, current_data: QuestState, setup_data: SetupData
    ) -> QuestState:
        """Return field definitions for the interaction form, based on current state."""
        ...

    def validate_field_change(
        self, field_id: str, new_value: Any, current_state: dict[str, Any]
    ) -> list[str]:
        """Validate a single field change before persisting. Returns list of errors."""
        ...
        return []

    def after_field_change(
        self, field_id: str, new_value: Any, current_state: dict[str, Any]
    ) -> list[dict[str, Any]]:
        """Called after a record is persisted. Returns list of side-effect actions."""
        ...

    def handle_event(
        self, event_name: str, current_state: dict[str, Any]
    ) -> list[dict[str, Any]]:
        """Handle system events (e.g. 'daily_reset'). Returns list of side-effect actions."""
        ...
