use std::error::Error;

use tauri::{AppHandle, Manager, State};
use mlua::LuaSerdeExt;

use crate::domain::QuestCreateAction;
use crate::managed::app_state::AppState;
use crate::repositories::PlanRepository;

// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/

#[tauri::command]
#[specta::specta]
pub fn debug_something(app: AppHandle) -> Result<(), String> {
    log::info!("Clicked the debug something.");

    let plans_dir = app
        .path()
        .resolve("plans", tauri::path::BaseDirectory::Resource)
        .map_err(|e| e.to_string())?;

    let plan_path = plans_dir.join("simple_quest");
    let main_path = plan_path.join("main.luau");

    log::info!("Loading plan from: {}", main_path.display());

    // --- 1. Создаём Luau VM ---
    let lua = mlua::Lua::new();

    // --- 2. Загружаем main.luau и получаем таблицу M ---
    let module: mlua::Table = lua
        .load(main_path.as_path())
        .set_name("main.luau")
        .eval()
        .map_err(|e| format!("Failed to load plan: {e}"))?;

    // --- 3. Достаём хук ---
    let on_create: mlua::Function = module
        .get("on_quest_create")
        .map_err(|e| format!("Hook 'on_quest_create' not found: {e}"))?;

    // --- 4. Готовим ctx и params ---
    let ctx = lua.create_table().map_err(|e| e.to_string())?;
    ctx.set("now", chrono::Utc::now().timestamp())
        .map_err(|e| e.to_string())?;

    let params = lua.create_table().map_err(|e| e.to_string())?;

    // --- 5. Вызываем хук ---
    let result_table: mlua::Table = on_create
        .call((ctx, params))
        .map_err(|e| format!("Hook failed: {e}"))?;

    // --- 6. Десериализуем Lua-таблицу в Rust-тип ---
    let action: QuestCreateAction = lua
        .from_value(mlua::Value::Table(result_table))
        .map_err(|e| format!("Failed to parse action: {e}"))?;

    log::info!("Parsed action: {action:?}");

    Ok(())
}
