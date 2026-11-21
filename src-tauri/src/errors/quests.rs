use serde::Serialize;
use thiserror::Error;

#[derive(specta::Type, Debug, Error, Serialize)]
#[serde(
    tag = "type",
    rename_all = "camelCase",
    rename_all_fields = "camelCase"
)]
pub enum QuestError {
    #[error("quest {quest_uuid} not found")]
    NotFound { quest_uuid: String },

    #[error("quest {quest_uuid} already exists")]
    AlreadyExists { quest_uuid: String },

    #[error("plan {plan_id} not found")]
    PlanNotFound { plan_id: String },
}