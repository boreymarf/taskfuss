use chrono::NaiveDateTime;
use diesel::prelude::*;

use crate::{db::plans, domain::PlanLanguage};

use super::permissions::PlanPermissions;

#[derive(Debug, Clone, Queryable, Selectable, Identifiable)]
#[diesel(table_name = plans)]
pub struct Plan {
    pub id: String,
    pub name: String,
    pub author: String,
    pub version: String,
    pub entry: String,
    #[diesel(deserialize_as = String)]
    pub language: PlanLanguage,
    #[diesel(deserialize_as = String)]
    pub permissions: PlanPermissions,
}

#[derive(Debug, Insertable, AsChangeset)]
#[diesel(table_name = plans)]
pub struct NewPlan {
    pub id: String,
    pub name: String,
    pub author: String,
    pub version: String,
    pub entry: String,
    #[diesel(serialize_as = String)]
    pub language: PlanLanguage,
    #[diesel(serialize_as = String)]
    pub permissions: PlanPermissions,
}
