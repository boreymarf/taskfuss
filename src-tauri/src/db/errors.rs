use thiserror::Error;

#[derive(Debug, Error, serde::Serialize)]
pub enum DbError {
    #[error("invalid UTF-8 in DB path: {0}")]
    NonUtf8Path(String),

    #[error("failed to create DB directory at {path}: {source}")]
    CreateDir {
        path: String,
        #[serde(skip)]
        #[source]
        source: std::io::Error,
    },

    #[error("failed to create a pool at {path}: {source}")]
    PoolOpen {
        path: String,
        #[serde(skip)]
        #[source]
        source: diesel::r2d2::PoolError,
    },

    #[error("failed to get connection from pool: {0}")]
    PoolConn(
        #[from]
        #[serde(skip)]
        diesel::r2d2::PoolError,
    ),

    #[error("failed to run migrations: {0}")]
    Migration(
        #[serde(skip)]
        #[source]
        Box<dyn std::error::Error + Send + Sync>,
    ),

    #[error("query failed: {0}")]
    Query(String),
}
