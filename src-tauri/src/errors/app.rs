use serde::Serialize;
use thiserror::Error;

use crate::{db::DbError, errors::{PlanRuntimeError, QuestError}};

use super::PlanError;

#[derive(specta::Type, Debug, Error, Serialize)]
#[serde(tag = "kind", content = "detail", rename_all = "camelCase")]
pub enum AppError {
    #[error(transparent)]
    Plan(#[from] PlanError),

    #[error(transparent)]
    Quest(#[from] QuestError),

    #[error(transparent)]
    Db(#[from] DbError),

    #[error(transparent)]
    PlanRuntime(#[from] PlanRuntimeError),

    // catch-all for anything that isn't a domain error
    #[error("internal error: {0}")]
    Internal(String),
}

impl From<diesel::result::Error> for AppError {
    fn from(e: diesel::result::Error) -> Self {
        DbError::Diesel(e).into()
    }
}
