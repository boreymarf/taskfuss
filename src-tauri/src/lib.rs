#![allow(warnings)]

mod commands;
mod db;
mod domain;
mod dto;
mod errors;
mod managed;
mod orchestrators;
mod repositories;
mod services;

use std::sync::Arc;

use colored::Colorize;
use tauri::Manager;
use tauri_plugin_log::{
    Target, TargetKind,
    fern::colors::{Color, ColoredLevelConfig},
};

use specta_typescript::Typescript;
use tauri_specta::{Builder, collect_commands};

// Learn more about Tauri commands at https://tauri.app/develop/calling-rust/
#[tauri::command]
#[specta::specta]
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
    // Spectra type generation
    let mut builder = Builder::new().commands(tauri_specta::collect_commands![
        greet,
        commands::plans::get_all_plans,
        commands::debug::debug_something,
        commands::plans::get_plan_settings_fields
    ]);

    #[cfg(debug_assertions)]
    builder
        .export(Typescript::default(), "../src/generated/bindings.ts")
        .expect("Failed to export typescript bindings");

    // Colors for logging
    let colors = ColoredLevelConfig::new()
        .error(Color::Red)
        .warn(Color::Yellow)
        .info(Color::Green)
        .debug(Color::TrueColor {
            r: 128,
            g: 128,
            b: 128,
        }) // Medium Grey
        .trace(Color::TrueColor {
            r: 80,
            g: 80,
            b: 80,
        }); // Darker Grey 

    let source_color = Color::Cyan;

    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(
            tauri_plugin_log::Builder::new()
                .clear_format()
                .targets([
                    Target::new(TargetKind::Stdout).format(move |out, message, record| {
                        let src = format!(
                            "{}:{}",
                            record.file().unwrap_or("?"),
                            record.line().unwrap_or(0),
                        );
                        out.finish(format_args!(
                            "{:<5} {} {}",
                            colors.color(record.level()),
                            src.color(source_color),
                            message
                        ))
                    }),
                    Target::new(TargetKind::LogDir { file_name: None }).format(
                        |out, message, record| {
                            let src = format!(
                                "{}:{}",
                                record.file().unwrap_or("?"),
                                record.line().unwrap_or(0),
                            );
                            out.finish(format_args!("{:<5} {} {}", record.level(), src, message))
                        },
                    ),
                    Target::new(TargetKind::Webview).format(|out, message, record| {
                        let src = format!(
                            "{}:{}",
                            record.file().unwrap_or("?"),
                            record.line().unwrap_or(0),
                        );
                        out.finish(format_args!("{:<5} {} {}", record.level(), src, message))
                    }),
                ])
                .level(log::LevelFilter::Debug)
                .build(),
        )
        .invoke_handler(builder.invoke_handler())
        .setup(move |app| {
            // Spectra events
            builder.mount_events(app);

            // DB Stuff
            let storage = Arc::new(managed::Storage::resolve(app)?);
            let db_path = storage.db_file();
            let pool = db::pool::create_pool(&db_path)?;
            let database = Arc::new(db::Database::new(pool));
            database.run_migrations()?;

            // Managed
            let plan_registry = Arc::new(managed::PlanRegistry::new());
            let plan_runtime = managed::PlanRuntime;

            // Repositories
            let quest_repo = Arc::new(repositories::QuestRepository::new());

            // Services
            let quest_service = Arc::new(services::QuestService::new(quest_repo.clone()));

            // Orchestration
            let quest_orchestrator =
                Arc::new(orchestrators::QuestOrchestrator::new(quest_service.clone()));

            // Setup
            let plans_dir = app
                .path()
                .resolve("plans", tauri::path::BaseDirectory::Resource)
                .map_err(|e| e.to_string())?;

            plan_registry
                .register_all_plans(&plans_dir)
                .map_err(|e| e.to_string())?;

            // Shove em inside
            let app_state = managed::AppState {
                storage,
                database,
                plan_registry,
                plan_runtime,
                quest_repo,
                quest_service,
                quest_orchestrator,
            };
            app.manage(app_state);

            log::info!("Tauri setup complete");

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            greet,
            get_app_data_dir,
            commands::debug::debug_something,
            commands::plans::get_plan_settings_fields,
            commands::plans::get_all_plans,
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
