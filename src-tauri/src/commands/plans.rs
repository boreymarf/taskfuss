use std::{collections::HashMap, path::Path};

use tauri::{AppHandle, Manager};

use crate::{domain::{Plan, Field}, errors::{AppError, PlanError}, managed::{AppState, PlanRuntime}};

/// A command that returns the list of plans currently registered.
///
/// # Errors
/// - Any plan error is returned as [`AppError`].
#[tauri::command]
#[specta::specta]
pub fn get_all_plans(app: AppHandle) -> Result<Vec<Plan>, AppError> {
    let app_state = app.state::<AppState>();

    Ok(app_state.plan_registry.get_all())
}

#[tauri::command]
#[specta::specta]
pub fn get_plan_settings_fields(
    app: AppHandle,
    plan_id: String,
) -> Result<HashMap<String, Field>, AppError> {
    let app_state = app.state::<AppState>();

    let plan = app_state
        .plan_registry
        .get(&plan_id)
        .ok_or_else(|| PlanError::NotFound { plan_id })?;

    let entry_path = Path::new(&plan.path).join(&plan.entry);

    let lua = PlanRuntime::create_lua_vm();
    let module = PlanRuntime::load_module(&lua, &entry_path)?;
    let ctx = PlanRuntime::create_context(&lua, &entry_path)?;

    PlanRuntime::get_settings_fields(&lua, &module, &entry_path, ctx)
}