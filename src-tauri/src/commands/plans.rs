use tauri::{AppHandle, Manager, State};

use crate::{db::errors::DbError, domain::Plan, errors::AppError, managed::app_state::AppState};

/// A command that returns the list of plans currently available in the database.
///
/// # Errors
/// - Any database error is passed through as [`DbError::Diesel`], wrapped in [`AppError`].
#[tauri::command]
#[specta::specta]
pub fn get_all_plans(app: AppHandle) -> Result<Vec<Plan>, AppError> {
    let app_state = app.state::<AppState>();
    let database = app_state.database.clone();
    let plan_orchestrator = app_state.plan_orchestrator.clone();

    database
        .transaction(|conn| plan_orchestrator.get_all_plans(conn))
        .inspect_err(|e| log::error!("get_all_plans failed: {e}"))
}
