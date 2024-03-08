CREATE TABLE "migration" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    migration_group_id INTEGER,
    desc VARCHAR(255) NOT NULL,
    executed_at TIMESTAMP NOT NULL
);