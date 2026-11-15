use std::{fs, path::Path, sync::Arc};

use log::{debug, info, warn};

use crate::{
    db::DbConn,
    domain::Plan,
    errors::{AppError, PlanError},
    services::PlanService,
};

pub struct PlanOrchestrator {
    plan_service: Arc<PlanService>,
}

impl PlanOrchestrator {
    pub fn new(plan_service: Arc<PlanService>) -> Self {
        Self { plan_service }
    }

    pub fn get_all_plans(&self, conn: &mut DbConn) -> Result<Vec<Plan>, AppError> {
        self.plan_service.get_all_plans(conn)
    }

    pub fn register_all_plans(&self, conn: &mut DbConn, plans_dir: &Path) -> Result<(), AppError> {
        info!("registering all plans from {}", plans_dir.display());

        // Clean up the table so we can reregister all plans
        let removed = self.plan_service.unregister_all_plans(conn)?;
        debug!("removed {} existing plans", removed);

        let dir_path = plans_dir.display().to_string();

        let entries = fs::read_dir(plans_dir).map_err(|source| PlanError::ReadPlansDir {
            path: dir_path.clone(),
            source,
        })?;

        let mut registered = 0usize;
        let mut failed = 0usize;

        for entry in entries {
            let entry = entry.map_err(|source| PlanError::ReadDirEntry {
                path: dir_path.clone(),
                source,
            })?;
            let path = entry.path();

            if !path.is_dir() {
                continue;
            }

            match self.plan_service.register_plan(conn, &path) {
                Ok(()) => {
                    registered += 1;
                    debug!("registered plan {}", path.display());
                }
                Err(e) => {
                    failed += 1;
                    warn!("failed to register plan {}: {}", path.display(), e);
                }
            }
        }

        info!(
            "finished registering plans: {} ok, {} failed",
            registered, failed
        );

        Ok(())
    }
}
