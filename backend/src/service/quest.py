# src/services/quest_service.py

import logging
from sqlalchemy.orm import Session

from src.classes.form_processor import FormProcessor
from src.domain.fields import FormFields
from src.domain.quest import Quest, QuestCreateRequest
from src.exceptions import SetupFormDataValidationError
from src.exceptions.generic import NotFoundError
from src.repositories.quest import QuestRepository
from src.service.plan import PlanService
from src.service.quest_settings import QuestSettingsService

logger = logging.getLogger(__name__)


class QuestService:
    def __init__(
        self,
        quest_repository: QuestRepository,
        plan_service: PlanService,
        quest_settings_service: QuestSettingsService,
    ):
        self.quest_repository = quest_repository
        self.plan_service = plan_service
        self.quest_settings_service = quest_settings_service

    def create_quest(
        self,
        db: Session,
        owner_id: int,
        data: QuestCreateRequest,
        *,
        auto_commit: bool = False,
    ) -> Quest:
        # Check if plan exists
        plan_db = self.plan_service.get_plan(db, data.plan_id)
        if not plan_db:
            raise NotFoundError("Plan", data.plan_id)

        # Get instance of the plan to work with
        plan_inst = self.plan_service.get_plan_instance(db, data.plan_id)

        # Get setup fields
        settings_fields: FormFields = plan_inst.get_settings_form()

        # Validate data
        form_processor = FormProcessor(settings_fields, data.settings)
        errors = form_processor.validate_data()
        if errors:
            raise SetupFormDataValidationError(data.plan_id, errors)

        # Plan should also validate data
        # TODO: Doesn't work yet since I don't know which error structs to make
        # errors = plan_inst.validate_settings_form(data.settings)

        # Create db records – now using injected services
        quest = self.quest_repository.add(db, owner_id, data.plan_id)
        _ = self.quest_settings_service.add(db, quest.id, data.settings)  # <-- changed

        if auto_commit:
            db.commit()

        return quest

    def get_all_by_user(self, db: Session, owner_id: int) -> list[Quest]:
        """Get all quests owned by a user."""
        return self.quest_repository.get_all_by_user(db, owner_id)
