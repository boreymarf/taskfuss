use serde::Serialize;
use thiserror::Error;

use crate::db::DbError;

use super::PlanError;

#[derive(Debug, Error, Serialize)]
#[serde(tag = "kind", content = "detail", rename_all = "camelCase")]
pub enum AppError {
    #[error(transparent)]
    Plan(#[from] PlanError),

    #[error(transparent)]
    Db(#[from] DbError),

    // catch-all for anything that isn't a domain error
    #[error("internal error: {0}")]
    Internal(String),
}

impl From<diesel::result::Error> for AppError {
    fn from(e: diesel::result::Error) -> Self {
        DbError::Diesel(e).into()
    }
}
