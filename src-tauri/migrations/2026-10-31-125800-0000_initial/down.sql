DROP TRIGGER IF EXISTS trg_records_updated_at;

DROP TRIGGER IF EXISTS trg_states_updated_at;
DROP TRIGGER IF EXISTS trg_quests_updated_at;

DROP INDEX IF EXISTS idx_records_latest;

DROP INDEX IF EXISTS idx_states_quest;
DROP INDEX IF EXISTS uq_states_open;

DROP INDEX IF EXISTS idx_quests_plan;

DROP TABLE IF EXISTS records;
DROP TABLE IF EXISTS quest_states;
DROP TABLE IF EXISTS quests;
DROP TABLE IF EXISTS plans;