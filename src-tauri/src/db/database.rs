use diesel_migrations::{EmbeddedMigrations, embed_migrations, MigrationHarness};

use crate::db::{DbConn, DbError, DbPool};

pub const MIGRATIONS: EmbeddedMigrations = embed_migrations!("migrations");

#[derive(Clone)]
pub struct Database {
    pool: DbPool,
}

impl Database {
    pub fn new(pool: DbPool) -> Self {
        Self { pool }
    }

    pub fn conn(&self) -> Result<DbConn, DbError> {
        Ok(self.pool.get()?)
    }

    pub fn run_migrations(&self) -> Result<(), DbError> {
        let mut conn = self.conn()?;
        conn.run_pending_migrations(MIGRATIONS)
            .map_err(|e| DbError::Migration(e))?;
        Ok(())
    }
}
