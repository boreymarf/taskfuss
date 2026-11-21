use std::sync::Arc;

use diesel::result::{DatabaseErrorKind, Error as DieselError};

use crate::db::DbConn;
use crate::domain::{NewQuest, Quest};
use crate::errors::{AppError, QuestError};
use crate::repositories::QuestRepository;

pub struct QuestService {
    quest_repo: Arc<QuestRepository>,
}

impl QuestService {
    pub fn new(quest_repo: Arc<QuestRepository>) -> Self {
        Self { quest_repo }
    }

    pub fn create(&self, conn: &mut DbConn, new_quest: NewQuest) -> Result<Quest, AppError> {
        // clone for errors
        let quest_uuid = new_quest.uuid.clone();
        let plan_id = new_quest.plan_id.clone();

        self.quest_repo
            .create(conn, new_quest)
            .map_err(|e| match e {
                DieselError::DatabaseError(DatabaseErrorKind::UniqueViolation, _) => {
                    QuestError::AlreadyExists { quest_uuid }.into()
                }
                DieselError::DatabaseError(DatabaseErrorKind::ForeignKeyViolation, _) => {
                    QuestError::PlanNotFound { plan_id }.into()
                }
                other => AppError::from(other),
            })
    }
}
