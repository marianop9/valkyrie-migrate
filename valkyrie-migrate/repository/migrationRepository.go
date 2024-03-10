package repository

import (
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/jmoiron/sqlx"
)

type MigrationRepo struct {
	db *sqlx.DB
}

func NewMigrationRepo(db *sqlx.DB) *MigrationRepo {
	return &MigrationRepo{
		db,
	}
}

func (repo *MigrationRepo) GetMigrations() ([]MigrationGroup, error) {
	query :=
		`SELECT mg.id,
		mg.name,
		mg.declared_date AS date,
		count(m.migration_group_id) AS migrationCount
	FROM migrationGroup mg	
		JOIN migration m on mg.id = m.migration_group_id
	GROUP BY m.migration_group_id`

	migrationGroups := []MigrationGroup{}

	if err := repo.db.Select(&migrationGroups, query); err != nil {
		return nil, err
	}

	return migrationGroups, nil
}

func (repo *MigrationRepo) ExecuteMigrations(migrations []*MigrationGroup) error {
	tx, _ := repo.db.Begin()
	defer tx.Rollback()

	for i := 0; i < len(migrations); i++ {
		fmt.Printf("executing group %s", migrations[i].Name)
		
		if err := repo.applyMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to execute group '%s', %v", migrations[i].Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *MigrationRepo) applyMigration(tx *sql.Tx, migration *MigrationGroup) error {
	for _, mig := range migration.Files {
		buf, err := io.ReadAll(mig)

		if err != nil {
			return err
		}

		if _, sqlErr := tx.Exec(string(buf)); sqlErr != nil {
			return fmt.Errorf("failed to execute %s: %v", migration.Name, sqlErr)
		}

		fmt.Printf("executed %s", migration.Name)
	}

	return nil
}

type MigrationGroup struct {
	Id             uint
	Name           string
	FolderName     string
	Date           time.Time
	Files          []io.Reader
	MigrationCount uint `db:"migrationCount"`
}

func (mg *MigrationGroup) AddFile(f io.Reader) {
	if mg.Files == nil {
		mg.Files = []io.Reader{
			f,
		}
	} else {
		mg.Files = append(mg.Files, f)
	}
}

type Migration struct {
	Name      string
	GroupName string
	FReader   io.Reader
}
