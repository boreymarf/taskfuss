use diesel::prelude::*;

use crate::db::schema::plan_snapshots;
use crate::db::{Database, DbError};
use crate::domain::{NewPlanSnapshot, PlanSnapshot};

pub(crate) struct PlanSnapshotRepository {
    db: Database,
}

impl PlanSnapshotRepository {
    pub(crate) fn new(db: Database) -> Self {
        Self { db }
    }

    pub(crate) fn insert(
        &self,
        plan_id: &str,
        path: &str,
        hash: &str,
        content: &str,
    ) -> Result<PlanSnapshot, DbError> {
        let mut conn = self.db.conn()?;

        let uuid = uuid::Uuid::now_v7().to_string();

        let new_snapshot = NewPlanSnapshot {
            uuid: &uuid,
            plan_id,
            path,
            hash,
            content,
        };

        diesel::insert_into(plan_snapshots::table)
            .values(&new_snapshot)
            .returning(PlanSnapshot::as_returning())
            .get_result(&mut conn)
            .map_err(DbError::Diesel)
    }

    pub(crate) fn get_all(&self) -> Result<Vec<PlanSnapshot>, DbError> {
        let mut conn = self.db.conn()?;

        plan_snapshots::table
            .order_by(plan_snapshots::created_at.desc())
            .load::<PlanSnapshot>(&mut conn)
            .map_err(DbError::Diesel)
    }
}

