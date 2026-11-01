use std::{error::Error, path::PathBuf};

use tauri::{Manager, Runtime};

pub(crate) struct Storage {
    data_dir: PathBuf,
    config_dir: PathBuf,
    cache_dir: PathBuf,
}

impl Storage {
    pub(crate) fn resolve<R: Runtime, M: Manager<R>>(m: &M) -> Result<Self, Box<dyn Error>> {
        let p = m.path();
        Ok(Self {
            data_dir: p.app_data_dir()?,
            config_dir: p.app_config_dir()?,
            cache_dir: p.app_cache_dir()?,
        })
    }

    pub(crate) fn db_file(&self) -> PathBuf {
        // Try to get env variable first if it exist and not empty
        if let Some(path) = std::env::var_os("TASKFUSS_DB_PATH").filter(|p| !p.is_empty()) {
            return PathBuf::from(path);
        }

        // Default
        self.data_dir.join("app.db")
    }
}
