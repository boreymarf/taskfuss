use std::fs;
use std::path::Path;

use diesel::r2d2::{ConnectionManager, Pool, PooledConnection};
use diesel::sqlite::SqliteConnection;

use crate::db::DbError;
use crate::db::Database;

// Docs: https://docs.rs/diesel/latest/diesel/r2d2/index.html

/// Type alias representing a database pool from which connections of type [`DbConn`] are taken.
/// Created via [`create_pool`]. Based exclusively on the Diesel crate's SQLite connection,
/// since other databases are not supported in Tauri.
pub type DbPool = Pool<ConnectionManager<SqliteConnection>>;

/// Type alias representing a database connection.
/// Can be obtained via [`DbPool::get`] or via the wrapper in [`Database::conn()`] as a shorthand.
pub type DbConn = PooledConnection<ConnectionManager<SqliteConnection>>;

/// Creates a pool for a SQLite database file at `db_path`.
/// Since it is SQLite, the database file is created automatically if it does not exist.
///
/// # Errors
///
/// - Returns [`DbError::CreateDir`] if creating parent directories fails.
/// - Returns [`DbError::NonUtf8Path`] if the path is not valid UTF-8.
/// - Returns [`DbError::PoolOpen`] if opening the pool fails.
pub fn create_pool(db_path: &Path) -> Result<DbPool, DbError> {
    if let Some(parent) = db_path.parent() {
        fs::create_dir_all(parent).map_err(|source| DbError::CreateDir {
            path: parent.display().to_string(),
            source,
        });
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
