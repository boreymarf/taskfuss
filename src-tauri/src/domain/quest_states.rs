use chrono::NaiveDateTime;
use diesel::prelude::*;

use crate::db::quest_states;

#[derive(serde::Serialize, specta::Type, Debug, Clone, Queryable, Selectable, Identifiable)]
#[diesel(table_name = quest_states)]
#[diesel(primary_key(uuid))]
pub struct QuestState {
    pub uuid: String,
    pub quest_uuid: String,
    pub internal: String,
    pub starts_at: Option<NaiveDateTime>,
    pub ends_at: Option<NaiveDateTime>,
}

#[derive(Debug, Insertable, AsChangeset)]
#[diesel(table_name = quest_states)]
pub struct NewQuestState {
    pub uuid: String,
    pub quest_uuid: String,
    pub internal: String,
    pub starts_at: Option<NaiveDateTime>,
    pub ends_at: Option<NaiveDateTime>,
}