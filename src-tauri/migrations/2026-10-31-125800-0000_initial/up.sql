CREATE TABLE plans (
    id          TEXT NOT NULL PRIMARY KEY,          -- unique plan id
    name        TEXT NOT NULL,
    author      TEXT NOT NULL,
    version     TEXT NOT NULL,
    entry_point TEXT NOT NULL,
    language    TEXT NOT NULL,
    permissions TEXT NOT NULL,                       -- JSON array/object

    CHECK (length(id) > 0),
    CHECK (language IN ('lua', 'luau')),
    CHECK (json_valid(permissions))
);

CREATE TABLE quests (
    uuid       TEXT      NOT NULL PRIMARY KEY,
    plan_id    TEXT      NOT NULL,
    title      TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S','now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S','now'))
);

CREATE INDEX idx_quests_plan_id ON quests (plan_id);

CREATE TRIGGER quests_set_updated_at
AFTER UPDATE ON quests
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE quests
    SET updated_at = strftime('%Y-%m-%d %H:%M:%S','now')
    WHERE uuid = OLD.uuid;
END;
