-- name: GetMigrations :many
SELECT mg.id,
    mg.name,
    count(m.migration_group_id) AS migrationCount
FROM migration_group mg	
    LEFT JOIN migration m on mg.id = m.migration_group_id
GROUP BY m.migration_group_id;

-- name: GetMigrationsByGroup :many
SELECT m.name,
    mg.name AS groupName
FROM migration m
    JOIN migration_group mg on mg.id = m.migration_group_id
WHERE m.migration_group_id = :id;

-- name: LogMigrationGroup :execresult
INSERT INTO migration_group (
    name
) VALUES (
    :name
)
RETURNING id;

-- name: LogMigration :exec
INSERT INTO migration (
    migration_group_id,
    name,
    executed_at
) VALUES (
    :groupId, :name, :executedAt
);