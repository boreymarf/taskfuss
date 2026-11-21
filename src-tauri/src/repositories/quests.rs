use diesel::prelude::*;
use log::trace;

use crate::db::{DbConn, quests};
use crate::domain::{NewQuest, Quest};

pub struct QuestRepository;

impl QuestRepository {
    pub fn new() -> Self {
        Self
    }

    pub fn create(&self, conn: &mut DbConn, new_quest: NewQuest) -> QueryResult<Quest> {
        diesel::insert_into(quests::table)
            .values(new_quest)
            .returning(Quest::as_returning())
            .get_result(conn)
            .inspect(|quest| trace!("inserted quest {}", quest.uuid))
    }

    pub fn get(&self, conn: &mut DbConn, quest_uuid: &str) -> QueryResult<Quest> {
        quests::table
            .find(quest_uuid)
            .first(conn)
            .inspect(|_| trace!("found quest {}", quest_uuid))
            .inspect_err(|e| {
                if matches!(e, diesel::result::Error::NotFound) {
                    trace!("quest {} not found", quest_uuid);
                }
            })
    }

    pub fn get_all(&self, conn: &mut DbConn) -> QueryResult<Vec<Quest>> {
        quests::table
            .load(conn)
            .inspect(|rows| trace!("loaded {} quests", rows.len()))
    }

    pub fn get_by_plan(&self, conn: &mut DbConn, plan_id: &str) -> QueryResult<Vec<Quest>> {
        quests::table
            .filter(quests::plan_id.eq(plan_id))
            .load(conn)
            .inspect(|rows| trace!("loaded {} quests for plan {}", rows.len(), plan_id))
    }

    pub fn delete(&self, conn: &mut DbConn, quest_uuid: &str) -> QueryResult<usize> {
        diesel::delete(quests::table.find(quest_uuid))
            .execute(conn)
            .inspect(|&n| trace!("deleted quest {} ({} rows)", quest_uuid, n))
    }
}