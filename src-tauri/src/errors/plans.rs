use serde::Serialize;
use thiserror::Error;

#[derive(Debug, Error, Serialize)]
#[serde(tag = "type", rename_all = "camelCase", rename_all_fields = "camelCase")]
pub enum PlanError {
    #[error("plan {plan_id} not found")]
    NotFound { plan_id: String },
}
