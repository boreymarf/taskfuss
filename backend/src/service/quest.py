import logging
from pprint import pprint

from sqlalchemy.orm import Session

from src.classes.form_processor import FormProcessor
from src.domain import (
    FormFields,
    Quest,
)
from src.domain.quest import QuestCreateRequest
from src.domain.quest_actions import CreateNewStateAction, QuestAction
from src.domain.quest_event import InitEvent
from src.exceptions import NotFoundError, SetupFormDataValidationError
from src.exceptions.generic import AlreadyExistsError
from src.exceptions.quest import UnknownQuestActionError
from src.plans.quest_context import QuestContext
from src.repositories.quest import QuestRepository
from src.repositories.quest_settings import QuestSettingsRepository
from src.repositories.quest_state import QuestStateRepository
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
    ) -> Quest:
        # Check if plan exists
        plan_db = PlanService.get_plan(db, data.plan_id)
        if not plan_db:
            raise NotFoundError("Plan", data.plan_id)

        # Get instance of the plan to work with
        plan_inst = PlanService.get_plan_instance(db, data.plan_id)

        # Get setup fields
        settings_fields: FormFields = plan_inst.get_settings_form()

        # Validate data
        form_processor = FormProcessor(settings_fields, data.settings)
        errors = form_processor.validate()
        if errors:
            raise SetupFormDataValidationError(data.plan_id, errors)

        # Plan should also validate data
        # TODO: Doesn't work yet since I don't know which error structs to make
        # errors = plan_inst.validate_settings_form(data.settings)

        # Create db records
        quest_db = QuestRepository.add(db, owner_id, data.plan_id)
        _ = QuestSettingsRepository.add(db, quest_db.id, data.settings)

        ctx = QuestContext(db, quest_db.id)

        # Get new state
        actions = plan_inst.handle_event(ctx, InitEvent())
        for action in actions:
            QuestService.handle_quest_action(db, action)

        if auto_commit:
            db.commit()

        return Quest.model_validate(quest_db)

    @staticmethod
    def handle_quest_action(db: Session, action: QuestAction):
        match action.discriminator:
            case "create_new_state":
                QuestService.handle_create_new_state_action(db, action)
            # case "update_settings":
            #     QuestService.handle_update_settings_action(db, action)
            # case "send_notification":
            #     QuestService.handle_send_notification_action(db, action)
            case _:
                logger.error(f"Unknown action type: {action.discriminator}")
                raise UnknownQuestActionError(action.discriminator)

    @staticmethod
    def handle_create_new_state_action(db: Session, action: CreateNewStateAction):

        quest_id = action.data.quest_id
        start_date = action.data.start_date
        end_date = action.data.end_date

        if QuestStateRepository.has_overlap(db, quest_id, start_date, end_date):
            logger.error(
                f"I can't create a new quest state, since there's already one! Quest id: '{quest_id}' for range ({start_date} to {end_date})"
            )
            raise AlreadyExistsError(
                entity_type="QuestState",
                details={
                    "quest_id": str(quest_id),
                    "start_date": start_date.isoformat() if start_date else None,
                    "end_date": end_date.isoformat() if end_date else None,
                    "message": "A state already exists for the given time range",
                },
            )

        quest_state = QuestStateRepository.add(db, action.data)

        pprint(quest_state)

    @staticmethod
    def get_all_by_user(db: Session, owner_id: int) -> list[Quest]:
        quests_db = QuestRepository.get_all_by_user(db, owner_id)

        quests = [Quest.model_validate(q) for q in quests_db]

        return quests
