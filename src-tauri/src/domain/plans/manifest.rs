use std::path::Path;

use serde::Deserialize;

use crate::errors::PlanError;

use super::NewPlan;

#[derive(Debug, Deserialize)]
pub struct RawManifest {
    schema_version: u32,
    plan: toml::Value,
}

pub trait IntoNewPlan {
    fn into_new_plan(self, path: &Path) -> Result<NewPlan, PlanError>;
}

#[derive(Debug, Deserialize, Clone)]
pub struct PlanSpecV1 {
    pub id: String,
    pub name: String,
    pub author: String,
    pub version: String,
    pub entry: String,
    pub language: String,
    #[serde(default)]
    pub permissions: Vec<String>,
}

#[derive(Debug, Clone)]
pub enum PlanSpec {
    V1(PlanSpecV1),
}

impl IntoNewPlan for PlanSpec {
    fn into_new_plan(self, path: &Path) -> Result<NewPlan, PlanError> {
        match self {
            PlanSpec::V1(v1) => v1.into_new_plan(path),
        }
    }
}

pub fn parse_manifest(path: &Path, contents: &str) -> Result<PlanSpec, PlanError> {
    let raw: RawManifest = toml::from_str(contents).map_err(|source| PlanError::ManifestParse {
        path: path.display().to_string(),
        source,
    })?;

    match raw.schema_version {
        1 => raw
            .plan
            .try_into()
            .map(PlanSpec::V1)
            .map_err(|source| PlanError::ManifestParse {
                path: path.display().to_string(),
                source,
            }),
        v => Err(PlanError::UnsupportedSchemaVersion {
            path: path.display().to_string(),
            version: v,
        }),
    }
}

impl IntoNewPlan for PlanSpecV1 {
    fn into_new_plan(self, path: &Path) -> Result<NewPlan, PlanError> {
        Ok(NewPlan {
            id: self.id,
            name: self.name,
            author: self.author,
            version: self.version,
            entry: self.entry,
            path: path.to_string_lossy().into_owned(),
            language: self
                .language
                .try_into()
                .map_err(PlanError::InvalidLanguage)?,
            permissions: self
                .permissions
                .try_into()
                .map_err(PlanError::InvalidPermissions)?,
        })
    }
}
