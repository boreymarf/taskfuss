use crate::{db::Database, managed::storage::Storage};

pub(crate) struct AppState {
    pub database: Database,
    pub storage: Storage
}

