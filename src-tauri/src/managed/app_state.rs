use std::sync::Arc;

use crate::db::Database;
use crate::managed::storage::Storage;
use crate::orchestrators::{PlanOrchestrator, QuestOrchestrator};
use crate::repositories::{PlanRepository, QuestRepository};
use crate::services::{PlanService, QuestService};

pub(crate) struct AppState {
    pub database: Arc<Database>,
    pub storage: Arc<Storage>,

    pub plan_repo: Arc<PlanRepository>,
    pub quest_repo: Arc<QuestRepository>,

    pub plan_service: Arc<PlanService>,
    pub quest_service: Arc<QuestService>,

    pub plan_orchestrator: Arc<PlanOrchestrator>,
    pub quest_orchestrator: Arc<QuestOrchestrator>,
}