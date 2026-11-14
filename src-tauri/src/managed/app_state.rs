use std::sync::Arc;

use crate::{db::Database, managed::storage::Storage, repositories::PlanRepository, services::PlanService};

pub(crate) struct AppState {
    pub database: Database,
    pub storage: Storage,

    pub plan_repo: Arc<PlanRepository>,

    pub plan_service: Arc<PlanService>,
}

