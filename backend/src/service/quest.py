from typing import Any

from sqlalchemy.orm import Session

from src.db import QuestDB
from src.dto.quest import QuestCreate
from src.exceptions import NotFoundError, PlanSetupValidationError
from src.plans.fields import FormFields
from src.plans.registry import PlanRegistry
from src.plans.setup_data import SetupData
from src.service.plan import PlanService


class QuestService:
    @staticmethod
    def _get_plan_instance(db: Session, plan_id: str):
        """Загружает план из БД и возвращает его экземпляр."""
        plan_registry = PlanService.get_plan(db, plan_id)
        if not plan_registry:
            raise NotFoundError("plan", plan_id)

        registry = PlanRegistry.model_validate(plan_registry)
        plan_cls = registry.import_class()
        return plan_cls()

    @staticmethod
    def validate_setup_form(
        db: Session, setup_form: FormFields, plan_id: str
    ) -> list[str]:
        instance = QuestService._get_plan_instance(db, plan_id)
        return instance.validate_setup_data(SetupData(form_data=setup_form))

    @staticmethod
    def create_quest(
        db: Session, owner_id: int, data: QuestCreate, *, auto_commit: bool = False
    ):
        # Валидация использует тот же метод получения инстанса
        issues = QuestService.validate_setup_form(db, data.setup_form, data.plan_id)
        if len(issues) != 0:
            raise PlanSetupValidationError(
                "Failed to validate a setup form", details={"issues": issues}
            )

        plan_instance = QuestService._get_plan_instance(db, data.plan_id)

        quest = QuestDB(
            owner_id=owner_id,
            setup_form_data=data.setup_form,
            plan_id=data.plan_id,
        )
        db.add(quest)
        db.flush()

        if auto_commit:
            db.commit()

        return quest
