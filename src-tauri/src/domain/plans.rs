use chrono::NaiveDateTime;
use diesel::prelude::*;

use crate::db::plans;

#[derive(Debug, Clone, Queryable, Selectable, Identifiable)]
#[diesel(table_name = plans)]
pub struct Plan {
    pub id: String,
    pub name: String,
    pub author: String,
    pub version: String,
    pub entry_point: String,
    pub language: String, // TODO: Later make it enum
    pub permissions: String,
}

#[derive(Debug, Insertable, AsChangeset)]
#[diesel(table_name = plans)]
pub struct NewPlan<'a> {
    pub id: &'a str,
    pub name: &'a str,
    pub author: &'a str,
    pub version: &'a str,
    pub entry_point: &'a str,
    pub language: &'a str, // TODO: Later make it enum
    pub permissions: &'a str,
}
