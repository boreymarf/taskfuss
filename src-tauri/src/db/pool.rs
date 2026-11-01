use std::fs;
use std::path::Path;

use diesel::r2d2::{ConnectionManager, Pool, PooledConnection};
use diesel::sqlite::SqliteConnection;

use crate::db::DbError;

// Docs: https://docs.rs/diesel/latest/diesel/r2d2/index.html

pub type DbPool = Pool<ConnectionManager<SqliteConnection>>;
pub type DbConn = PooledConnection<ConnectionManager<SqliteConnection>>;

pub fn create_pool(db_path: &Path) -> Result<DbPool, DbError> {
    if let Some(parent) = db_path.parent() {
        fs::create_dir_all(parent).map_err(|source| DbError::CreateDir {
            path: parent.display().to_string(),
            source,
        })?;
    }

    let db_path_str = db_path
        .to_str()
        .ok_or_else(|| DbError::NonUtf8Path(db_path.display().to_string()))?;

    let manager = ConnectionManager::<SqliteConnection>::new(db_path_str);

    let pool = Pool::builder()
        .event_handler(Box::new(super::pool_logger::DbPoolLogger))
        .build(manager)
        .map_err(|source| DbError::PoolOpen {
            path: db_path.display().to_string(),
            source,
        })?;

    log::info!("Pool was created for the path {db_path_str}");

    Ok(pool)
}
