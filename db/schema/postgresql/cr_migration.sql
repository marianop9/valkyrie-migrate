CREATE TABLE migration (
    id SERIAL,
    migration_group_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    executed_at TIMESTAMP NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_migration
        FOREIGN KEY (migration_group_id) REFERENCES "migration_group" (id)
);
