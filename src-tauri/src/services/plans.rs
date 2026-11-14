use std::path::Path;
use std::sync::Arc;

use crate::repositories::PlanRepository;
use crate::errors::PlanError;


pub struct PlanService {
    plan_repo: Arc<PlanRepository>,
}

impl PlanService {
    pub fn new(plan_repo: Arc<PlanRepository>) -> Self {
        Self { plan_repo }
    }

    pub fn register_plan(&self, path: &Path) -> Result<(), PlanError> {

        // The plan path must be a directory
        if !path.is_dir() {
            return Err(PlanError::NotADirectory {
                path: path.to_string_lossy().to_string(),
            });
        }

        let manifest_path = path.join("manifest.toml");
        if !manifest_path.is_file() {
            return Err(PlanError::ManifestNotFound {
                path: manifest_path.to_string_lossy().to_string(),
            });
        }
        
        let manifest_contents = std::fs::read_to_string(&manifest_path).map_err(|e| PlanError::ManifestRead {
            path: manifest_path.to_string_lossy().to_string(),
            source: e,
        })?;

        let manifest: toml::Value = toml::from_str(&manifest_contents).map_err(|e| PlanError::ManifestParse {
            path: manifest_path.to_string_lossy().to_string(),
            source: e,
        })?;

        dbg!(&manifest);

        Ok(())
    }
}