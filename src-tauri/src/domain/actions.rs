use chrono::NaiveDateTime;
use serde::{Deserialize, Serialize};

use crate::domain::PlanField;

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
#[serde(tag = "type")]
pub enum QuestCreateAction {
    #[serde(rename = "create_quest_state")]
    CreateQuestState(CreateQuestStateActionData),
}

#[derive(Debug, Clone, Serialize, Deserialize, specta::Type)]
pub struct CreateQuestStateActionData {
    pub internal: serde_json::Value,
    pub starts_at: Option<i64>,
    pub ends_at: Option<i64>,
    #[serde(default)]
    pub fields: Vec<PlanField>,
}