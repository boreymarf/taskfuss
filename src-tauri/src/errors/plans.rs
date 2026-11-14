use serde::Serialize;
use thiserror::Error;

#[derive(Debug, Error, Serialize)]
#[serde(tag = "type", rename_all = "camelCase", rename_all_fields = "camelCase")]
pub enum PlanError {
    #[error("plan {plan_id} not found")]
    NotFound { plan_id: String },

    #[error("path {path} is not a directory")]
    NotADirectory { path: String },

    #[error("manifest not found at {path}")]
    ManifestNotFound { path: String },

    #[error("failed to read manifest at {path}")]
    ManifestRead {
        path: String,
        #[source]
        #[serde(skip_serializing)]
        source: std::io::Error,
    },

    #[error("failed to parse manifest at {path}")]
    ManifestParse {
        path: String,
        #[source]
        #[serde(skip_serializing)]
        source: toml::de::Error,
    },
}
