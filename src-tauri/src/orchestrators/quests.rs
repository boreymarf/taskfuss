use std::{fs, path::Path, sync::Arc};

use log::{debug, info, warn};

use crate::db::DbConn;
use crate::domain::{NewQuest, Quest};
use crate::errors::{AppError, QuestError};
use crate::services::QuestService;

pub struct QuestOrchestrator {
    quest_service: Arc<QuestService>,
}

impl QuestOrchestrator {
    pub fn new(quest_service: Arc<QuestService>) -> Self {
        Self { quest_service }
    }

    pub fn create_quest(quest: NewQuest) {}
}
