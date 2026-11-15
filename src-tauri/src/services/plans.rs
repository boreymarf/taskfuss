use diesel::result::{DatabaseErrorKind, Error as DieselError};
use std::path::Path;
use std::sync::Arc;

use crate::db::{DbConn, DbError};
use crate::domain::{IntoNewPlan, NewPlan, Plan, RawManifest, parse_manifest};
use crate::errors::{AppError, PlanError};
use crate::repositories::PlanRepository;

pub struct PlanService {
    plan_repo: Arc<PlanRepository>,
}

impl PlanService {
    pub fn new(plan_repo: Arc<PlanRepository>) -> Self {
        Self { plan_repo }
    }

    pub fn get_all_plans(&self, conn: &mut DbConn) -> Result<Vec<Plan>, AppError> {
        self.plan_repo.get_all(conn).map_err(|e| DbError::Diesel(e).into())
    }

    pub fn register_plan(&self, conn: &mut DbConn, path: &Path) -> Result<(), AppError> {
        if !path.is_dir() {
            return Err(PlanError::NotADirectory {
                path: path.display().to_string(),
            }
            .into());
        }

        let manifest_path = path.join("manifest.toml");
        if !manifest_path.is_file() {
            return Err(PlanError::ManifestNotFound {
                path: manifest_path.display().to_string(),
            }
            .into());
        }

        let manifest_contents =
            std::fs::read_to_string(&manifest_path).map_err(|e| PlanError::ManifestRead {
                path: manifest_path.display().to_string(),
                source: e,
            })?;

        let manifest = parse_manifest(&manifest_path, &manifest_contents)?;
        let new_plan = manifest.into_new_plan()?;
        let plan_id = new_plan.id.clone(); // so we can put it in the error message

        self.plan_repo.create(conn, new_plan).map_err(|e| match e {
            DieselError::DatabaseError(DatabaseErrorKind::UniqueViolation, _) => {
                PlanError::AlreadyExists { plan_id }.into()
            }
            other => AppError::Db(DbError::Diesel(other)),
        })?;

        Ok(())
    }

    pub fn unregister_all_plans(&self, conn: &mut DbConn) -> Result<usize, AppError> {
        self.plan_repo
            .delete_all(conn)
            .map_err(|e| DbError::Diesel(e).into())
    }
}
