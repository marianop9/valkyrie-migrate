package sqliteRepo

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/marianop9/valkyrie-migrate/internal/models"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
	queries "github.com/marianop9/valkyrie-migrate/internal/repository/queries/sqlite"
)

type SqliteRepo struct {
	db      *sql.DB
	queries *queries.Queries
}

func NewMigrationRepo(db *sql.DB) *SqliteRepo {
	return &SqliteRepo{
		db:      db,
		queries: queries.New(db),
	}
}

func (repo *SqliteRepo) EnsureCreated() error {
	return repository.EnsureCreated(repo.db)
}

func (repo *SqliteRepo) GetMigrations() ([]models.MigrationGroup, error) {
	queryRows, err := repo.queries.GetMigrations(context.TODO())
	if err != nil {
		return nil, err
	}

	migrationGroups := migGroupFromQueryList(queryRows)

	for i := range migrationGroups {
		group := &migrationGroups[i]

		migs, err := repo.getMigrationsByGroup(group.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to read migrations from group '%s': %v", group.Name, err)
		}
		group.Migrations = migs
	}

	return migrationGroups, nil
}

func (repo *SqliteRepo) getMigrationsByGroup(groupId uint) ([]models.Migration, error) {
	queryResult, err := repo.queries.GetMigrationsByGroup(context.TODO(), int64(groupId))
	if err != nil {
		return nil, err
	}

	return migFromQueryList(queryResult), nil
}

func (repo *SqliteRepo) ExecuteMigrations(migrations []*models.MigrationGroup) error {
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

		fmt.Printf("done executing group %s\n", migrations[i].Name)
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
			GroupId: int64(group.Id),
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
		MigrationCount: int(queryRow.Migrationcount),
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
