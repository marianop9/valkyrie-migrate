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

func (repo *MigrationRepo) EnsureCreated() error {
	return EnsureCreated(repo.db)
}

func (repo *MigrationRepo) GetMigrations() ([]MigrationGroup, error) {
	query := `SELECT mg.id,
		mg.name,
		count(m.migration_group_id) AS migrationCount
	FROM migrationGroup mg	
		LEFT JOIN migration m on mg.id = m.migration_group_id
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
		fmt.Printf("executing group %s:\n", migrations[i].Name)

		if err := repo.applyMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to execute group '%s', %v", migrations[i].Name, err)
		}

		if err := logMigration(tx, migrations[i]); err != nil {
			return fmt.Errorf("failed to log group '%s', %v", migrations[i].Name, err)
		}

		fmt.Printf("done executing group %s:\n", migrations[i].Name)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *MigrationRepo) applyMigration(tx *sql.Tx, migration *MigrationGroup) error {
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

func logMigration(tx *sql.Tx, group *MigrationGroup) error {
	groupCmd := `INSERT INTO migrationGroup (
		name
	) VALUES (
		$1
	)
	RETURNING id`

	result, err := tx.Exec(groupCmd, group.Name)
	if err != nil {
		return err
	}

	groupId, err := result.LastInsertId()
	if err != nil {
		return err
	} 
	group.Id = uint(groupId)

	migrationCmd := `INSERT INTO migration (
		migration_group_id,
		name,
		executed_at
	) VALUES (
		$1, $2, $3
	)`

	logTime := time.Now()
	for _, mig := range group.Migrations {
		_, err := tx.Exec(migrationCmd,
			group.Id,
			mig.Name,
			logTime)
		
		if err != nil {
			return err
		}
	}

	return nil
}


type MigrationGroup struct {
	Id             uint
	Name           string
	Files          []io.Reader
	Migrations     []Migration
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

func (mg *MigrationGroup) AddMigration(mig Migration) {
	if mg.Migrations == nil {
		mg.Migrations = []Migration{
			mig,
		}
	} else {
		mg.Migrations = append(mg.Migrations, mig)
	}
}

type Migration struct {
	Name      string
	GroupName string
	FReader   io.Reader
}
