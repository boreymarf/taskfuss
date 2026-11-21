use chrono::NaiveDateTime;
use diesel::prelude::*;
use serde::{Deserialize, Serialize};
use specta::Type;

use crate::db::quests;

#[derive(
    Serialize,
    Deserialize,
    Type,
    Debug,
    Clone,
    Queryable,
    Selectable,
    Identifiable,
)]
#[diesel(table_name = quests)]
#[diesel(primary_key(uuid))]
pub struct Quest {
    pub uuid: String,
    pub plan_id: String,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Insertable)] 
#[diesel(table_name = quests)] 
pub struct NewQuest { 
    pub uuid: String, 
    pub plan_id: String 
}