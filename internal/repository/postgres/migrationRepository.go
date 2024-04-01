package postgresRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/marianop9/valkyrie-migrate/internal/models"
	queries "github.com/marianop9/valkyrie-migrate/internal/repository/queries/postgresql"
)

var ErrInconsistenMigrationSchema = errors.New("found only one of the required tables: migration or migration_group")

type MigrationRepo struct {
	db      *sql.DB
	queries *queries.Queries
}

func NewMigrationRepo(db *sql.DB) *MigrationRepo {
	return &MigrationRepo{
		db:      db,
		queries: queries.New(db),
	}
}

func (repo *MigrationRepo) EnsureCreated() error {
	migrationTables := []string{
		"migration_group",
		"migration",
	}

	query := `SELECT count(1)
		FROM information_schema.tables 
		WHERE table_schema = 'public'
			AND table_name IN ($1, $2);`

	var tableCount int
	err := repo.db.QueryRow(query, migrationTables[0], migrationTables[1]).Scan(&tableCount)
	if err != nil {
		return err
	}

	if tableCount == len(migrationTables) {
		fmt.Println("migrations tables exist")
		return nil
	} else if tableCount != 0 {
		return ErrInconsistenMigrationSchema
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// create migration schema
	if err = createMigrationTables(tx); err != nil {
		return err
	}

	return tx.Commit()
}

func createMigrationTables(tx *sql.Tx) error {
	fmt.Println("creating table 'migration_group'...")

	cmd1 := `CREATE TABLE "migration_group" (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL
	);`

	if _, err := tx.Exec(cmd1); err != nil {
		return err
	}

	fmt.Println(`creating table 'migration'...`)

	cmd2 := `CREATE TABLE migration (
		id SERIAL,
		migration_group_id INTEGER NOT NULL,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT fk_migration FOREIGN KEY (migration_group_id) REFERENCES "migration_group" (id)
	);`

	if _, err := tx.Exec(cmd2); err != nil {
		return err
	}

	return nil
}

func (repo *MigrationRepo) GetMigrations() ([]models.MigrationGroup, error) {
	queryRows, err := repo.queries.GetMigrations(context.TODO())
	if err != nil {
		return nil, err
	}

	migrationGroups := migGroupFromQueryList(queryRows)

	for i := 0; i < len(migrationGroups); i++ {
		group := migrationGroups[i]

		migs, err := repo.getMigrationsByGroup(group.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to read migrations from group '%s': %v", group.Name, err)
		}
		group.Migrations = migs
	}

	return migrationGroups, nil
}

func (repo *MigrationRepo) getMigrationsByGroup(groupId uint) ([]models.Migration, error) {
	queryResult, err := repo.queries.GetMigrationsByGroup(context.TODO(), int32(groupId))
	if err != nil {
		return nil, err
	}

	return migFromQueryList(queryResult), nil
}

func (repo *MigrationRepo) ExecuteMigrations(migrations []*models.MigrationGroup) error {
	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := 0; i < len(migrations); i++ {
		fmt.Printf("executing group %s:\n", migrations[i].Name)

		if err := applyMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to execute group '%s', %v", migrations[i].Name, err)
		}

		if err := logMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to log group '%s', %v", migrations[i].Name, err)
		}

		fmt.Printf("done executing group %s:\n", migrations[i].Name)
	}

	return tx.Commit()
}

func applyMigration(tx *sql.Tx, migration *models.MigrationGroup) error {
	for _, mig := range migration.Migrations {
		buf, err := io.ReadAll(mig.FReader)

		if err != nil {
			return err
		}

		if _, sqlErr := tx.Exec(string(buf)); sqlErr != nil {
			return fmt.Errorf("failed to execute %s: %v", migration.Name, sqlErr)
		}

		fmt.Printf("\t * executed %s\n", mig.Name)
	}

	return nil
}

func logMigration(tx *sql.Tx, group *models.MigrationGroup) error {
	// logs are executed manualy because pgx doesn't support returning LastInsertId when
	// executing a query, so QueryRow is used instead.
	// sqlc is prob overkill for this project anyways :)

	migrationGroupCmd := `INSERT INTO migration_group (
		name
	) VALUES (
		$1
	)
	RETURNING id;`

	var groupId uint

	err := tx.QueryRow(migrationGroupCmd, group.Name).Scan(&groupId)
	if err != nil {
		return err
	}
	group.Id = groupId

	migrationCmd := `INSERT INTO migration (
		migration_group_id,
		name,
		executed_at
	) VALUES ($1, $2, $3);`

	logTime := time.Now()
	for _, mig := range group.Migrations {
		if _, err := tx.Exec(migrationCmd, groupId, mig.Name, logTime); err != nil {
			return err
		}
	}

	return nil
}

func migGroupFromQuery(queryRow *queries.GetMigrationsRow) models.MigrationGroup {
	return models.MigrationGroup{
		Id:             uint(queryRow.ID),
		Name:           queryRow.Name,
		MigrationCount: int(queryRow.MigrationCount),
	}
}

func migGroupFromQueryList(queryResult []queries.GetMigrationsRow) []models.MigrationGroup {
	list := make([]models.MigrationGroup, len(queryResult))

	for i := 0; i < len(queryResult); i++ {
		list[i] = migGroupFromQuery(&queryResult[i])
	}

	return list
}

func migFromQuery(row *queries.GetMigrationsByGroupRow) models.Migration {
	return models.Migration{
		Name:      row.Name,
		GroupName: row.Groupname,
	}
}

func migFromQueryList(queryResult []queries.GetMigrationsByGroupRow) []models.Migration {
	list := make([]models.Migration, len(queryResult))

	for i := 0; i < len(queryResult); i++ {
		list[i] = migFromQuery(&queryResult[i])
	}

	return list
}
