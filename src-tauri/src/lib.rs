mod db;
mod managed;
mod models;

use tauri::Manager;

// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
#[tauri::command]
fn greet(name: &str) -> String {
    format!("Hello, {}! You've been greeted from Rust!", name)
}

#[tauri::command]
fn get_app_data_dir(app: tauri::AppHandle) -> Result<String, String> {
    app.path()
        .app_data_dir()
        .map(|p| p.to_string_lossy().into_owned())
        .map_err(|e| e.to_string())
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .setup(|app| {
            let storage = managed::storage::Storage::resolve(app)?;
            let db_path = storage.db_file();
            let pool = db::pool::create_pool(&db_path)?;
            let database = db::Database::new(pool);

            // Run migrations
            database.run_migrations()?;

            // Shove em inside
            let app_state = managed::app_state::AppState { storage, database };
            app.manage(app_state);

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![greet, get_app_data_dir])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
