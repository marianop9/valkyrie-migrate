package postgresRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
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

	query := `SELECT table_name
		FROM information_schema.tables 
		WHERE table_schema = 'public'
			AND table_name IN ($1, $2);`

	rows, err := repo.db.Query(query, migrationTables[0], migrationTables[1])
	if err != nil {
		return err
	}

	defer rows.Close()

	foundTables := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return err
		}

		foundTables = append(foundTables, name)
	}

	if len(foundTables) == len(migrationTables) {
		fmt.Println("migrations tables exist")
		return nil
	} else if len(foundTables) != 0 {
		return ErrInconsistenMigrationSchema
	}

	// create migration_group table
	if err = createMigrationGroupTable(repo.db); err != nil {
		return err
	}

	// create migration table
	if err = createMigrationTable(repo.db); err != nil {
		return err
	}

	return nil
}

func createMigrationGroupTable(db *sql.DB) error {
	fmt.Println("creating table 'migration_group'...")

	buf, err := os.ReadFile("./db/schema/postgresql/cr_migrationGroup.sql")
	if err != nil {
		return err
	}

	if _, sqlErr := db.Exec(string(buf)); sqlErr != nil {
		return sqlErr
	}

	return nil
}

func createMigrationTable(db *sql.DB) error {
	fmt.Println(`creating table 'migration'...`)

	buf, err := os.ReadFile("./db/schema/postgresql/cr_migration.sql")
	if err != nil {
		return err
	}

	if _, sqlErr := db.Exec(string(buf)); sqlErr != nil {
		return sqlErr
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

	txQuery := repo.queries.WithTx(tx)
	
	for i := 0; i < len(migrations); i++ {
		fmt.Printf("executing group %s:\n", migrations[i].Name)

		if err := applyMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to execute group '%s', %v", migrations[i].Name, err)
		}

		if err := logMigration(txQuery, migrations[i]); err != nil {
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

func logMigration(tx *queries.Queries, group *models.MigrationGroup) error {
	result, err := tx.LogMigrationGroup(context.TODO(), group.Name)
	if err != nil {
		return err
	}

	if groupId, err := result.LastInsertId(); err != nil {
		return err
	} else {
		group.Id = uint(groupId)
	}

	logTime := time.Now()
	for _, mig := range group.Migrations {
		migrationParams := queries.LogMigrationParams{
			MigrationGroupID: int32(group.Id),
			Name: mig.Name,
			ExecutedAt: logTime,
		}

		if err := tx.LogMigration(context.TODO(), migrationParams); err != nil {
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
