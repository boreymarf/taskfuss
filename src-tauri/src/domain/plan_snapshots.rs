use chrono::NaiveDateTime;
use diesel::prelude::*;
use serde::Serialize;

use crate::db::schema::plan_snapshots;

#[derive(Debug, Clone, Queryable, Selectable, Identifiable, Serialize)]
#[diesel(table_name = plan_snapshots)]
#[diesel(primary_key(uuid))]
pub struct PlanSnapshot {
    pub uuid: String,
    /// Internal id from the plan manifest
    pub plan_id: String,
    /// Path to a plan (to watch and stuff)
    pub path: String,
    pub hash: String,
    pub content: String,
    pub created_at: NaiveDateTime,
}

#[derive(Debug, Insertable)]
#[diesel(table_name = plan_snapshots)]
pub struct NewPlanSnapshot<'a> {
    pub uuid: &'a str,
    pub plan_id: &'a str,
    pub path: &'a str,
    pub hash: &'a str,
    pub content: &'a str,
}
