package app

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/marianop9/valkyrie-migrate/app/helpers"
	"github.com/marianop9/valkyrie-migrate/app/migrations"
	"github.com/marianop9/valkyrie-migrate/app/repository"
)

type MigrateApp struct {
	repo *repository.MigrationRepo
}

func NewMigrateApp(repo *repository.MigrationRepo) *MigrateApp {
	return &MigrateApp{
		repo,
	}
}

func (app MigrateApp) Run(dsn string) error {
	folderNames := []string{
		"./migrations",
		"../migrations",
	}

	var baseFolderName string

	// get migrations folder
	var dirEntries []os.DirEntry

	for _, fname := range folderNames {
		entries, err := os.ReadDir(fname)

		if err == nil && len(entries) > 0 {
			dirEntries = entries
			baseFolderName = fname
			break
		}
	}

	fmt.Println("files: ", len(dirEntries))
	if len(dirEntries) == 0 {
		return fmt.Errorf("no migrations found in folders %+v", folderNames)
	}

	if err := app.repo.EnsureCreated(); err != nil {
		fmt.Println("failed to create migration tables")
		return err
	}

	// retrieve migrations from folder
	migrationGroups, err := migrations.GetMigrationGroups(dirEntries)

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
		if existingMigFolder := helpers.FindByName(existingMigrations, migrationFolder.Name); existingMigFolder == nil {
			// new migration folder
			migrationGroupsToApply = append(migrationGroupsToApply, migrationFolder)
		}
	}

	if len(migrationGroupsToApply) == 0 {
		fmt.Println("database is up to date. Exiting...")
		return nil
	}

	for i, newFolder := range migrationGroupsToApply {
		// read dir to get the files
		migrationFolderPath := path.Join(baseFolderName, newFolder.Name)
		migrationFiles, err := os.ReadDir(migrationFolderPath)

		if err != nil {
			return fmt.Errorf("failed to read migration folder: %v", err)
		}
		
		for _, migrationFile := range migrationFiles {
			if fReader, err := os.Open(path.Join(migrationFolderPath, migrationFile.Name())); err != nil {
				return errors.Join(fmt.Errorf("failed to read file %v", migrationFile.Name()), err)
			} else {
				newMig := repository.Migration{
					Name:      migrationFile.Name(),
					GroupName: migrationGroupsToApply[i].Name,
					FReader:   fReader,
				}
				migrationGroupsToApply[i].AddMigration(newMig)
			}
		}
	}

	if err := app.repo.ExecuteMigrations(migrationGroupsToApply); err != nil {
		return err
	}

	return nil
}
