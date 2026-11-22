use std::collections::HashMap;
use std::fs;
use std::path::Path;
use std::sync::RwLock;

use log::{debug, info, trace, warn};

use crate::domain::{IntoNewPlan, NewPlan, Plan, parse_manifest};
use crate::errors::PlanError;

pub struct PlanRegistry {
    plans: RwLock<HashMap<String, Plan>>,
}

impl PlanRegistry {
    pub fn new() -> Self {
        Self {
            plans: RwLock::new(HashMap::new()),
        }
    }

    pub fn register_all_plans(&self, plans_dir: &Path) -> Result<(), PlanError> {
        info!("registering all plans from {}", plans_dir.display());

        let removed = self.unregister_all_plans();
        debug!("removed {} existing plans", removed);

        let dir_path = plans_dir.display().to_string();

        let entries = fs::read_dir(plans_dir).map_err(|source| PlanError::ReadPlansDir {
            path: dir_path.clone(),
            source,
        })?;

        let mut registered = 0usize;
        let mut failed = 0usize;

        for entry in entries {
            let entry = entry.map_err(|source| PlanError::ReadDirEntry {
                path: dir_path.clone(),
                source,
            })?;
            let path = entry.path();

            if !path.is_dir() {
                continue;
            }

            match self.register_plan(&path) {
                Ok(_) => {
                    registered += 1;
                    debug!("registered plan {}", path.display());
                }
                Err(e) => {
                    failed += 1;
                    warn!("failed to register plan {}: {}", path.display(), e);
                }
            }
        }

        info!(
            "finished registering plans: {} ok, {} failed",
            registered, failed
        );

        Ok(())
    }

    pub fn register_plan(&self, path: &Path) -> Result<Plan, PlanError> {
        if !path.is_dir() {
            return Err(PlanError::NotADirectory {
                path: path.display().to_string(),
            });
        }

        let manifest_path = path.join("manifest.toml");
        if !manifest_path.is_file() {
            return Err(PlanError::ManifestNotFound {
                path: manifest_path.display().to_string(),
            });
        }

        let manifest_contents =
            fs::read_to_string(&manifest_path).map_err(|e| PlanError::ManifestRead {
                path: manifest_path.display().to_string(),
                source: e,
            })?;

        let manifest = parse_manifest(&manifest_path, &manifest_contents)?;
        let new_plan = manifest.into_new_plan(path)?;

        self.create(new_plan)
    }

    pub fn create(&self, new_plan: NewPlan) -> Result<Plan, PlanError> {
        let mut plans = self.plans.write().expect("plans lock poisoned");

        if plans.contains_key(&new_plan.id) {
            return Err(PlanError::AlreadyExists {
                plan_id: new_plan.id,
            });
        }

        let plan = Plan {
            id: new_plan.id,
            name: new_plan.name,
            author: new_plan.author,
            version: new_plan.version,
            entry: new_plan.entry,
            path: new_plan.path,
            language: new_plan.language,
            permissions: new_plan.permissions,
        };

        trace!("inserted plan {}", plan.id);
        plans.insert(plan.id.clone(), plan.clone());
        Ok(plan)
    }

    pub fn get(&self, plan_id: &str) -> Option<Plan> {
        let plans = self.plans.read().expect("plans lock poisoned");
        let plan = plans.get(plan_id).cloned();
        match &plan {
            Some(_) => trace!("found plan {}", plan_id),
            None => trace!("plan {} not found", plan_id),
        }
        plan
    }

    pub fn get_all(&self) -> Vec<Plan> {
        let plans = self.plans.read().expect("plans lock poisoned");
        let rows: Vec<Plan> = plans.values().cloned().collect();
        trace!("loaded {} plans", rows.len());
        rows
    }

    pub fn update(&self, plan_id: &str, changes: NewPlan) -> Option<Plan> {
        let mut plans = self.plans.write().expect("plans lock poisoned");

        if !plans.contains_key(plan_id) {
            trace!("plan {} not found for update", plan_id);
            return None;
        }

        let plan = Plan {
            id: plan_id.to_string(),
            name: changes.name,
            author: changes.author,
            version: changes.version,
            entry: changes.entry,
            path: changes.path,
            language: changes.language,
            permissions: changes.permissions,
        };

        trace!("updated plan {}", plan_id);
        plans.insert(plan_id.to_string(), plan.clone());
        Some(plan)
    }

    pub fn delete(&self, plan_id: &str) -> bool {
        let mut plans = self.plans.write().expect("plans lock poisoned");
        let removed = plans.remove(plan_id).is_some();
        trace!("deleted plan {} ({})", plan_id, removed);
        removed
    }

    pub fn unregister_all_plans(&self) -> usize {
        let mut plans = self.plans.write().expect("plans lock poisoned");
        let n = plans.len();
        plans.clear();
        trace!("deleted all plans ({})", n);
        n
    }
}