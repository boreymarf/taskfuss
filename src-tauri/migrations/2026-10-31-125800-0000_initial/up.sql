-- quests
CREATE TABLE quests (
    uuid               TEXT NOT NULL PRIMARY KEY,
    plan_snapshot_uuid TEXT NOT NULL,
    title              TEXT,
    created_at         TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
    updated_at         TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),

    FOREIGN KEY (plan_snapshot_uuid)
        REFERENCES plan_snapshots (uuid)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX idx_quests_plan_snapshot_uuid
    ON quests (plan_snapshot_uuid);

-- update time when updating the quest
CREATE TRIGGER quests_set_updated_at
AFTER UPDATE ON quests
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE quests
    SET updated_at = strftime('%Y-%m-%d %H:%M:%S', 'now')
    WHERE uuid = OLD.uuid;
END;

-- plans
CREATE TABLE plan_snapshots (
    uuid        TEXT NOT NULL PRIMARY KEY,
    plan_id     TEXT NOT NULL,
    path        TEXT NOT NULL,
    hash        TEXT NOT NULL,
    content     TEXT NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),

    CHECK (length(plan_id) > 0),
    CHECK (length(hash)    > 0),
    CHECK (json_valid(content))
);

-- can't have same id + hash combo
-- guards against copies
CREATE UNIQUE INDEX idx_plan_snapshots_plan_id_hash
    ON plan_snapshots (plan_id, hash);

-- gimme the most recent snapshot
CREATE INDEX idx_plan_snapshots_plan_id_created_at
    ON plan_snapshots (plan_id, created_at DESC);
