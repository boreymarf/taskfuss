CREATE TABLE quests (
    id          INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    title       TEXT,
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now')),
    updated_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%S', 'now'))
);

CREATE TRIGGER quests_set_updated_at
AFTER UPDATE ON quests
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
BEGIN
    UPDATE quests
    SET updated_at = strftime('%Y-%m-%d %H:%M:%S', 'now')
    WHERE id = OLD.id;
END;
