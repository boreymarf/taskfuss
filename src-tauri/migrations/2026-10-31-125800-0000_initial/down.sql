DROP TRIGGER IF EXISTS quests_set_updated_at;
DROP TABLE   IF EXISTS quests;

DROP INDEX IF EXISTS idx_plan_snapshots_plan_id_hash;
DROP INDEX IF EXISTS idx_plan_snapshots_plan_id_created_at;
DROP TABLE IF EXISTS plan_snapshots;
