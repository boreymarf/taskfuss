use std::collections::HashMap;
use std::path::Path;
use std::sync::Arc;

use mlua::{Function, Lua, LuaSerdeExt, Table, Value};

use crate::domain::{HookName, Field, QuestCreateAction, QuestCreateParams};
use crate::errors::{AppError, PlanRuntimeError};

pub struct PlanRuntime;

impl PlanRuntime {
    pub fn create_lua_vm() -> Lua {
        Lua::new()
    }

    pub fn load_module(lua: &Lua, path: &Path) -> Result<Table, AppError> {
        let module: Table =
            lua.load(path)
                .set_name("main")
                .eval()
                .map_err(|e| PlanRuntimeError::LoadModule {
                    path: path.to_string_lossy().into_owned(),
                    source: e,
                })?;

        Ok(module)
    }

    pub fn create_context(lua: &Lua, path: &Path) -> Result<Table, AppError> {
        let err = |e| PlanRuntimeError::ContextCreate {
            path: path.to_string_lossy().into_owned(),
            source: e,
        };

        let ctx = lua.create_table().map_err(err)?;
        ctx.set("now", chrono::Utc::now().timestamp())
            .map_err(err)?;

        Ok(ctx)
    }

    pub fn get_settings_fields(
        lua: &Lua,
        module: &Table,
        path: &Path,
        ctx: Table,
    ) -> Result<HashMap<String, Field>, AppError> {
        let hook = HookName::GetSettingsFields;
        let path_str = path.to_string_lossy().into_owned();

        let get_fields: Function =
            module
                .get(hook.as_str())
                .map_err(|e| PlanRuntimeError::HookResolve {
                    path: path_str.clone(),
                    hook,
                    source: e,
                })?;

        let result_value: Value =
            get_fields
                .call(ctx)
                .map_err(|e| PlanRuntimeError::HookInvoke {
                    path: path_str.clone(),
                    hook,
                    source: e,
                })?;

        let fields: HashMap<String, Field> =
            lua.from_value(result_value)
                .map_err(|e| PlanRuntimeError::HookDeserialize {
                    path: path_str,
                    hook,
                    source: e,
                })?;

        Ok(fields)
    }

    pub fn on_quest_create(
        lua: &Lua,
        module: &Table,
        path: &Path,
        ctx: Table,
        params: QuestCreateParams,
    ) -> Result<Vec<QuestCreateAction>, AppError> {
        let hook = HookName::OnQuestCreate;
        let path_str = path.to_string_lossy().into_owned();

        let params_value: Value =
            lua.to_value(&params)
                .map_err(|e| PlanRuntimeError::HookInvoke {
                    path: path_str.clone(),
                    hook,
                    source: e,
                })?;

        let on_create: Function =
            module
                .get(hook.as_str())
                .map_err(|e| PlanRuntimeError::HookResolve {
                    path: path_str.clone(),
                    hook,
                    source: e,
                })?;

        let result_value: Value =
            on_create
                .call((ctx, params_value))
                .map_err(|e| PlanRuntimeError::HookInvoke {
                    path: path_str.clone(),
                    hook,
                    source: e,
                })?;

        let actions: Vec<QuestCreateAction> =
            lua.from_value(result_value)
                .map_err(|e| PlanRuntimeError::HookDeserialize {
                    path: path_str,
                    hook,
                    source: e,
                })?;

        Ok(actions)
    }
}
