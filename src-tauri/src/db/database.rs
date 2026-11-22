use diesel::Connection;
use diesel_migrations::{EmbeddedMigrations, MigrationHarness, embed_migrations};

use crate::db::{DbConn, DbError, DbPool};
use crate::managed::AppState;
use std::sync::Arc;

pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

/// State struct stored in [`AppState`] and serves as the interface for working with the database.
/// Holds the database pool.
/// 
/// Should be cloned and stored in the app state using [`Arc`].
#[derive(Clone)]
pub struct Database {
    pool: DbPool,
}

impl Database {
    /// Creates a new `Database` from a `DbPool`.
    /// The `DbPool` is created separately via the [`create_pool`] function.
    pub fn new(pool: DbPool) -> Self {
        Self { pool }
    }

    /// Creates a connection from the pool.
    pub fn conn(&self) -> Result<DbConn, DbError> {
        Ok(self.pool.get()?)
    }

    /// Runs all available Diesel migrations.
    /// Returns [`DbError::Migration`] on any error.
    pub fn run_migrations(&self) -> Result<(), DbError> {
        let mut conn = self.conn()?;
        conn.run_pending_migrations(MIGRATIONS)
            .map_err(DbError::Migration)?;
        Ok(())
    }

    /// Creates a connection `conn` and runs function `F` within a transaction.
    ///
    /// # Example
    ///
    /// ```no_run
    /// # use crate::db::{Database, DbError};
    /// # fn example(db: &Database) -> Result<(), DbError> {
    /// db.transaction(|conn| {
    ///     // perform operations with `conn`
    ///     Ok(())
    /// })?;
    /// # Ok(())
    /// # }
    /// ```
    pub fn transaction<T, E, F>(&self, f: F) -> Result<T, E>
    where
        F: FnOnce(&mut DbConn) -> Result<T, E>,
        E: From<DbError> + From<diesel::result::Error>,
    {
        let mut conn = self.conn()?;
        conn.transaction(f)
    }
}
