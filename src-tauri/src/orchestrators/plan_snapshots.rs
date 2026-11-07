use crate::db::DbError;
use crate::domain::PlanSnapshot;
use crate::services::PlanSnapshotService;

pub(crate) struct PlanSnapshotOrchestrator {
    service: PlanSnapshotService,
}

impl PlanSnapshotOrchestrator {
    pub(crate) fn new(service: PlanSnapshotService) -> Self {
        Self { service }
    }

    pub(crate) fn collect_all(&self) -> Result<Vec<PlanSnapshot>, DbError> {
        self.service.get_all()
    }
}
