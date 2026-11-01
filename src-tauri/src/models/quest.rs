use chrono::NaiveDateTime;
use diesel::prelude::*;
use serde::{Deserialize, Serialize};

use crate::db::schema::quests;

#[derive(Debug, Clone, Queryable, Selectable, Identifiable, Serialize)]
#[diesel(table_name = quests)]
pub struct Quest {
    pub id: i32,
    pub title: Option<String>,
    pub created_at: NaiveDateTime,
    pub updated_at: NaiveDateTime,
}

#[derive(Debug, Insertable)]
#[diesel(table_name = quests)]
pub struct NewQuest<'a> {
    pub title: Option<&'a str>,
}

#[derive(Debug, AsChangeset, Deserialize)]
#[diesel(table_name = quests)]
pub struct UpdateQuest {
    pub title: Option<String>,
}
