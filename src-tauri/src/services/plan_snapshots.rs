use tauri::AppHandle;

use crate::db::DbError;
use crate::domain::PlanSnapshot;
use crate::repositories::PlanSnapshotRepository;

pub(crate) struct PlanSnapshotService {
    repo: PlanSnapshotRepository,
}

impl PlanSnapshotService {
    pub(crate) fn new(repo: PlanSnapshotRepository) -> Self {
        Self { repo }
    }

    pub(crate) fn insert(
        &self,
        plan_id: &str,
        path: &str,
        hash: &str,
        content: &str,
    ) -> Result<PlanSnapshot, DbError> {
        self.repo.insert(plan_id, path, hash, content)
    }

    pub(crate) fn get_all(&self) -> Result<Vec<PlanSnapshot>, DbError> {
        self.repo.get_all()
    }
}
