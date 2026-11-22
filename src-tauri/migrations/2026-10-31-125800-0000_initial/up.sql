CREATE TABLE plans (
    id          TEXT NOT NULL PRIMARY KEY,          -- unique plan id
    name        TEXT NOT NULL,
    author      TEXT NOT NULL,
    version     TEXT NOT NULL,
    entry       TEXT NOT NULL,
    path        TEXT NOT NULL,
    language    TEXT NOT NULL,
    permissions TEXT NOT NULL,                       -- JSON array/object

    CHECK (length(id) > 0),
    CHECK (language IN ('lua', 'luau')),
    CHECK (json_valid(permissions))
);

CREATE TABLE quests (
    uuid       TEXT      NOT NULL PRIMARY KEY,
    plan_id    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
);
CREATE INDEX idx_quests_plan ON quests(plan_id);

CREATE TRIGGER trg_quests_updated_at
AFTER UPDATE OF uuid, plan_id ON quests
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE quests
       SET updated_at = strftime('%Y-%m-%d %H:%M:%f','now')
     WHERE uuid = NEW.uuid;
END;

CREATE TABLE quest_states (
    uuid       TEXT      NOT NULL PRIMARY KEY,
    quest_uuid TEXT      NOT NULL REFERENCES quests(uuid) ON DELETE CASCADE,
    internal   TEXT      NOT NULL DEFAULT '{}',
    starts_at  TIMESTAMP,
    ends_at    TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now')),
    updated_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
);
CREATE INDEX idx_states_quest ON quest_states(quest_uuid, starts_at, ends_at);

CREATE TRIGGER trg_states_updated_at
AFTER UPDATE OF quest_uuid, internal, starts_at, ends_at ON quest_states
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE quest_states
       SET updated_at = strftime('%Y-%m-%d %H:%M:%f','now')
     WHERE uuid = NEW.uuid;
END;

CREATE TABLE records (
    uuid       TEXT      NOT NULL PRIMARY KEY,
    state_uuid TEXT      NOT NULL REFERENCES quest_states(uuid) ON DELETE CASCADE,
    field_id   TEXT      NOT NULL,
    author     TEXT      NOT NULL CHECK (author IN ('user','plan')),
    data       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f','now'))
);
CREATE INDEX idx_records_latest ON records(state_uuid, field_id, created_at DESC);