from pprint import pprint
from typing import Any

from sqlalchemy.orm import Session

from src.db import QuestDB
from src.domain import FormFields, QuestSettingsCreate, validate_form
from src.dto.quest import QuestCreate
from src.exceptions import NotFoundError, SetupFormDataValidationError
from src.service.plan import PlanService


class QuestService:
    @staticmethod
    def create_quest(
        db: Session, owner_id: int, data: QuestCreate, *, auto_commit: bool = False
    ):
        # Check if plan exists
        plan_db = PlanService.get_plan(db, data.plan_id)
        if not plan_db:
            raise NotFoundError("Plan", data.plan_id)

        # Get instance of the plan to work with
        plan_inst = PlanService.get_plan_instance(db, data.plan_id)

        # Get setup fields
        setup_fields: FormFields = plan_inst.get_setup_fields()

        # Validate data
        errors = validate_form(setup_fields, data.setup_form_data)
        if errors:
            raise SetupFormDataValidationError(data.plan_id, errors)

        # Collect setup data
        # Can later have other info like date
        setup_data = QuestSettingsCreate(form_data=setup_fields)

        # Plan should also validate data
        # TODO: Doesn't work yet since I don't know which error structs to make
        # errors = plan_inst.validate_setup_data(setup_data)

    @staticmethod
    def add_quest(
        db: Session, data: QuestCreate, *, auto_commit: bool = False
    ) -> QuestDB:
        quest = QuestDB(owner_id=data.owner_id, plan_id=data.plan_id)
        db.add(quest)
        db.flush()
        db.refresh(quest)
        if auto_commit:
            db.commit()
        return quest

    @staticmethod
    def add_quest_settings(
        db: Session, data: QuestSettingsCreate, *, auto_commit: bool = False
    ) -> QuestSettingsDB:
        settings_db = QuestSettingsDB(
            quest_id=data.quest_id,
            form_data=data.form_data
        )
        db.add(settings_db)
        db.flush()
        db.refresh(settings_db)
        if auto_commit:
            db.commit()
        return settings_db

    @staticmethod
    def add_quest_state(
        db: Session, data: QuestStateCreate, *, auto_commit: bool = False
    ) -> QuestStateDB:
        # data содержит все поля, включая quest_id
        state_db = QuestStateDB(**data.model_dump())
        db.add(state_db)
        db.flush()
        db.refresh(state_db)
        if auto_commit:
            db.commit()
        return state_db



    # @staticmethod
    # def _get_plan_instance(db: Session, plan_id: str):
    #     """Загружает план из БД и возвращает его экземпляр."""
    #     plan_registry = PlanService.get_plan(db, plan_id)
    #     if not plan_registry:
    #         raise NotFoundError("plan", plan_id)
    #
    #     registry = PlanRegistry.model_validate(plan_registry)
    #     plan_cls = registry.import_class()
    #     return plan_cls()
    #
    # @staticmethod
    # def validate_setup_form(
    #     db: Session, setup_form: FormFields, plan_id: str
    # ) -> list[str]:
    #     instance = QuestService._get_plan_instance(db, plan_id)
    #     return instance.validate_setup_data(SetupData(form_data=setup_form))
    #
    # @staticmethod
    # def create_quest(
    #     db: Session, owner_id: int, data: QuestCreate, *, auto_commit: bool = False
    # ):
    #     # Валидация использует тот же метод получения инстанса
    #     issues = QuestService.validate_setup_form(db, data.setup_form, data.plan_id)
    #     if len(issues) != 0:
    #         raise PlanSetupValidationError(
    #             "Failed to validate a setup form", details={"issues": issues}
    #         )
    #
    #     plan_instance = QuestService._get_plan_instance(db, data.plan_id)
    #
    #     quest = QuestDB(
    #         owner_id=owner_id,
    #         setup_form_data=data.setup_form,
    #         plan_id=data.plan_id,
    #     )
    #     db.add(quest)
    #     db.flush()
    #
    #     if auto_commit:
    #         db.commit()
    #
    #     return quest
