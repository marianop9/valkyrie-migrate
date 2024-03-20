package valkyrie

import (
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/migrations"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
)

type MigrateApp struct {
	repo *repository.MigrationRepo
}

func NewMigrateApp(repo *repository.MigrationRepo) *MigrateApp {
	return &MigrateApp{
		repo,
	}
}

// Creates a new migration instance connected to the specified database
func NewMigration(db *sql.DB, dbDriver string) *MigrateApp {
	repo := repository.NewMigrationRepo(db)
	return NewMigrateApp(repo)
}

func (app MigrateApp) Run(migrationFolder string) error {
	// get migrations directory
	dirEntries, err := os.ReadDir(migrationFolder)

	if err != nil {
		return err
	}

	if len(dirEntries) == 0 {
		return fmt.Errorf("no migrations found in folder: %+v", migrationFolder)
	}
	fmt.Println("found groups: ", len(dirEntries))

	if err := checkMigrationSubfolders(dirEntries); err != nil {
		return err
	}

	if err := app.repo.EnsureCreated(); err != nil {
		fmt.Println("failed to create migration tables")
		return err
	}

	// retrieve migrations from folder
	migrationGroups, err := migrations.GetMigrationGroups(migrationFolder, dirEntries)

	if err != nil {
		return err
	} else if len(migrationGroups) == 0 {
		fmt.Println("no migration groups found")
		return nil
	}

	// retrieve db migrations
	existingMigrations, err := app.repo.GetMigrations()

	if err != nil {
		return errors.Join(errors.New("failed to retrieve migrations from db"), err)
	}

	// find differences
	migrationGroupsToApply := make([]*repository.MigrationGroup, 0)
	for _, migrationFolder := range migrationGroups {
		if existingMigFolder := helpers.FindMigrationGroup(existingMigrations, migrationFolder.Name); existingMigFolder == nil {
			// new migration group
			migrationGroupsToApply = append(migrationGroupsToApply, migrationFolder)
		} else if existingMigFolder.MigrationCount != migrationFolder.MigrationCount {
			// migration group has new migrations to apply
			migrationsToApply := make([]repository.Migration, 0)

			for _, migrationFile := range migrationFolder.Migrations {
				if existingMigFile := helpers.FindMigration(existingMigFolder.Migrations, migrationFile.Name); existingMigFile != nil {
					migrationsToApply = append(migrationsToApply, migrationFile)
				}
			}

			groupToApply := &repository.MigrationGroup{
				Name:           migrationFolder.Name,
				Migrations:     migrationsToApply,
				MigrationCount: len(migrationsToApply),
			}
			migrationGroupsToApply = append(migrationGroupsToApply, groupToApply)
		}
	}

	if len(migrationGroupsToApply) == 0 {
		fmt.Println("database is up to date. Exiting...")
		return nil
	}

	fmt.Println("Groups to execute:")
	for _, group := range migrationGroupsToApply {
		fmt.Printf("* %s\n", group.Name)
		fmt.Printf("\t - migrations: %v\n", group.MigrationCount)
	}
	fmt.Printf("********\n\n")

	for _, groupToApply := range migrationGroupsToApply {
		// get the handles for files we need to migrate
		migrationFolderPath := path.Join(migrationFolder, groupToApply.Name)

		for i := 0; i < len(groupToApply.Migrations); i++ {
			migration := &groupToApply.Migrations[i]

			fReader, err := os.Open(path.Join(migrationFolderPath, migration.Name))
			if err != nil {
				return errors.Join(fmt.Errorf("failed to read file %v", migration.Name), err)
			}
			migration.FReader = fReader
		}
	}

	if err := app.repo.ExecuteMigrations(migrationGroupsToApply); err != nil {
		return err
	}

	return nil
}

func checkMigrationSubfolders(migrationFolderEntries []fs.DirEntry) error {
	isNotDir := func(dir os.DirEntry) bool {
		return !dir.IsDir()
	}

	if helpers.Any(migrationFolderEntries, isNotDir) {
		return fmt.Errorf("migrations folder may only contain subfolders representing migration groups")
	}

	return nil
}
