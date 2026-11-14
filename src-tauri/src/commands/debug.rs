use std::error::Error;

use tauri::{AppHandle, Manager, State};

use crate::repositories::PlanRepository;
use crate::managed::app_state::AppState;

// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

#[tauri::command]
pub fn debug_something(app: AppHandle) -> Result<(), String> {
    log::info!("Clicked the debug something.");

    let plans_dir = app
        .path()
        .resolve("plans", tauri::path::BaseDirectory::Resource)
        .map_err(|e| e.to_string())?;

    log::info!(
        "Plans dir: {} (exists: {})",
        plans_dir.display(),
        plans_dir.exists()
    );

    let app_state = app.state::<AppState>();
    
    let plan_service = app_state.plan_service.clone();
    let plan_path = plans_dir.join("simple_quest");

    plan_service.register_plan(&plan_path);

    Ok(())
}