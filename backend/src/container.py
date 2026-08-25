from dependency_injector import containers, providers

from src.repositories.plan import PlanRepository
from src.repositories.quest import QuestRepository
from src.repositories.quest_settings import QuestSettingsRepository
from src.repositories.quest_state import QuestStateRepository
from src.repositories.record import RecordRepository
from src.repositories.user import UserRepository
from src.service.plan import PlanService
from src.service.quest import QuestService
from src.service.quest_settings import QuestSettingsService
from src.service.quest_state import QuestStateService
from src.service.record import RecordService
from src.service.user import UserService


class Container(containers.DeclarativeContainer):
    wiring_config = containers.WiringConfiguration(
        modules=[
            "src.api.routers.user",
            "src.api.routers.quest",
            "src.api.routers.plan",
            "src.api.routers.quest_state",
            "src.api.routers.record",
            "src.api.routers.config",
        ]
    )

    quest_repository = providers.Singleton(QuestRepository)
    plan_repository = providers.Singleton(PlanRepository)
    user_repository = providers.Singleton(UserRepository)
    quest_state_repository = providers.Singleton(QuestStateRepository)
    quest_settings_repository = providers.Singleton(QuestSettingsRepository)
    record_repository = providers.Singleton(RecordRepository)

    user_service = providers.Singleton(
        UserService,
        user_repository=user_repository,
    )

    plan_service = providers.Singleton(
        PlanService,
        plan_repository=plan_repository,
    )

    quest_state_service = providers.Singleton(
        QuestStateService,
        quest_state_repository=quest_state_repository,
    )
    quest_settings_service = providers.Singleton(
        QuestSettingsService,
        quest_settings_repository=quest_settings_repository,
    )

    quest_service = providers.Singleton(
        QuestService,
        quest_repository=quest_repository,
        plan_service=plan_service,
        quest_settings_service=quest_settings_service,   
    )

    record_service = providers.Singleton(
        RecordService,
        record_repository=record_repository,
        quest_repository=quest_repository,
        quest_state_repository=quest_state_repository,
    )
