use serde::Serialize;
use thiserror::Error;

use crate::domain::HookName;

#[derive(specta::Type, Debug, Error, Serialize)]
#[serde(
    tag = "type",
    rename_all = "camelCase",
    rename_all_fields = "camelCase"
)]
pub enum PlanRuntimeError {
    #[error("failed to create runtime context for a plan {path}: {source}")]
    ContextCreate {
        path: String,
        #[source]
        #[serde(skip)]
        #[specta(skip)]
        source: mlua::prelude::LuaError,
    },

    #[error("runtime failed to load a module {path}: {source}")]
    LoadModule {
        path: String,
        #[source]
        #[serde(skip)]
        #[specta(skip)]
        source: mlua::prelude::LuaError,
    },

    #[error("hook `{hook}` not exported by plan {path}: {source}")]
    HookResolve {
        path: String,
        hook: HookName,
        #[source]
        #[serde(skip)]
        #[specta(skip)]
        source: mlua::prelude::LuaError,
    },

    #[error("hook `{hook}` threw in plan {path}: {source}")]
    HookInvoke {
        path: String,
        hook: HookName,
        #[source]
        #[serde(skip)]
        #[specta(skip)]
        source: mlua::prelude::LuaError,
    },

    #[error("hook `{hook}` returned an invalid payload in plan {path}: {source}")]
    HookDeserialize {
        path: String,
        hook: HookName,
        #[source]
        #[serde(skip)]
        #[specta(skip)]
        source: mlua::prelude::LuaError,
    },
}
