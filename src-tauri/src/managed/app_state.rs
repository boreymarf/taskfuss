use std::sync::Arc;

use crate::db::Database;
use crate::managed::storage::Storage;
use crate::orchestrators::plans::PlanOrchestrator;
use crate::repositories::PlanRepository;
use crate::services::PlanService;

pub(crate) struct AppState {
    pub database: Arc<Database>,
    pub storage: Arc<Storage>,

    pub plan_repo: Arc<PlanRepository>,

    pub plan_service: Arc<PlanService>,

    pub plan_orchestrator: Arc<PlanOrchestrator>,
}
