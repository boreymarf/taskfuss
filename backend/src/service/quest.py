import logging
from pprint import pprint
from typing import Any

from sqlalchemy.orm import Session

from src.db import QuestDB
from src.domain import FormFields, QuestCreate, QuestSettingsCreate, validate_form
from src.domain.quest import QuestCreateRequest
from src.domain.quest_actions import CreateNewStateAction, QuestAction, QuestActionBase
from src.domain.quest_event import InitEvent
from src.exceptions import NotFoundError, SetupFormDataValidationError
from src.exceptions.quest import UnknownQuestActionError
from src.repositories.quest import QuestRepository
from src.repositories.quest_settings import QuestSettingsRepository
from src.service.plan import PlanService

logger = logging.getLogger(__name__)


class QuestService:
    @staticmethod
    def create_quest(
        db: Session,
        owner_id: int,
        data: QuestCreateRequest,
        *,
        auto_commit: bool = False,
    ):
        # Check if plan exists
        plan_db = PlanService.get_plan(db, data.plan_id)
        if not plan_db:
            raise NotFoundError("Plan", data.plan_id)

        # Get instance of the plan to work with
        plan_inst = PlanService.get_plan_instance(db, data.plan_id)

        # Get setup fields
        settings_fields: FormFields = plan_inst.get_settings_form()

        # Validate data
        errors = validate_form(settings_fields, data.settings)
        if errors:
            raise SetupFormDataValidationError(data.plan_id, errors)

        # Collect setup data
        # Can later have other info like date
        # setup_data = QuestSettingsCreate(form_data=setup_fields)

        # Plan should also validate data
        # TODO: Doesn't work yet since I don't know which error structs to make
        # errors = plan_inst.validate_settings_form(data.settings)

        # Create db records
        quest_db = QuestRepository.add(db, owner_id, data.plan_id)
        quest_settings_db = QuestSettingsRepository.add(db, quest_db.id, data.settings)

        # Get new state
        actions = plan_inst.handle_event(InitEvent())
        for action in actions:
            QuestService.handle_quest_action(db, action)

    @staticmethod
    def handle_quest_action(db: Session, action: QuestAction):
        match action.discriminator:
            # case "create_new_state":
            #     QuestService.handle_create_new_state_action(db, action)
            # case "update_settings":
            #     QuestService.handle_update_settings_action(db, action)
            # case "send_notification":
            #     QuestService.handle_send_notification_action(db, action)
            case _:
                logger.error(f"Unknown action type: {action.discriminator}")
                raise UnknownQuestActionError(action.discriminator)

    @staticmethod
    def handle_create_new_state_action(db: Session, action: CreateNewStateAction):
        pass
