use std::path::Path;

use diesel::r2d2::{ConnectionManager, Pool, PooledConnection};
use diesel::sqlite::SqliteConnection;

use crate::db::DbError;

// Docs: https://docs.rs/diesel/latest/diesel/r2d2/index.html

pub type DbPool = Pool<ConnectionManager<SqliteConnection>>;
pub type DbConn = PooledConnection<ConnectionManager<SqliteConnection>>;

pub fn create_pool(db_path: &Path) -> Result<DbPool, DbError> {
    let db_path_str = db_path
        .to_str()
        .ok_or_else(|| DbError::NonUtf8Path(db_path.display().to_string()))?;

    let manager = ConnectionManager::<SqliteConnection>::new(db_path_str);

    let pool = Pool::builder()
        .build(manager)
        .map_err(|source| DbError::PoolOpen {
            path: db_path.display().to_string(),
            source,
        })?;

    Ok(pool)
}
