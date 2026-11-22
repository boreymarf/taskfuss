use std::sync::Arc;

use crate::db::Database;
use crate::managed::{PlanRegistry, PlanRuntime, Storage};
use crate::orchestrators::QuestOrchestrator;
use crate::repositories::QuestRepository;
use crate::services::QuestService;

pub struct AppState {
    pub database: Arc<Database>,
    pub storage: Arc<Storage>,
    pub plan_registry: Arc<PlanRegistry>,
    pub plan_runtime: PlanRuntime,

    pub quest_repo: Arc<QuestRepository>,

    pub quest_service: Arc<QuestService>,

    pub quest_orchestrator: Arc<QuestOrchestrator>,
}