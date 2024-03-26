package migrations

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/models"
)

const (
	dateFmt = "20060102"
)

func GetMigrationGroups(migrationDir string, dirEntries []os.DirEntry) ([]*models.MigrationGroup, error) {
	migrationGroups := make([]*models.MigrationGroup, 0)

	for _, dir := range dirEntries {
		dirName := dir.Name()

		files, err := os.ReadDir(path.Join(migrationDir, dirName))
		if err != nil {
			return nil, fmt.Errorf("failed to read dir '%s' - %v", dirName, err)
		}

		group := models.MigrationGroup{
			Name:           dirName,
			MigrationCount: len(files),
			Migrations:     []models.Migration{},
		}

		if err := checkFileExtension(files, dirName); err != nil {
			return nil, err
		}

		for _, file := range files {
			fileName := file.Name()
			fileNameParts := strings.Split(fileName, "_")

			if len(fileNameParts) < 2 {
				return nil, fmt.Errorf("(%s): file name doesn't match expected format (yyyymmdd_description)", fileName)
			}

			_, err := time.Parse(dateFmt, fileNameParts[0]) // yyyymmdd
			if err != nil {
				return nil, fmt.Errorf("(%s): file date doesn't match expected format (yyyymmdd)", fileName)
			}

			migration := models.Migration{
				Name:      fileName,
				GroupName: group.Name,
			}

			group.Migrations = append(group.Migrations, migration)
		}

		migrationGroups = append(migrationGroups, &group)
	}

	return migrationGroups, nil
}

func checkFileExtension(migrationGroupFiles []fs.DirEntry, folderName string) error {
	isDir := func(entry os.DirEntry) bool {
		return entry.IsDir()
	}

	if helpers.Any(migrationGroupFiles, isDir) {
		return fmt.Errorf("migration group folder may not contain nested subfolders. (%s)", folderName)
	}

	isNotSql := func(file os.DirEntry) bool {
		return path.Ext(file.Name()) != ".sql"
	}

	if helpers.Any(migrationGroupFiles, isNotSql) {
		return fmt.Errorf("migration group folder may only contain sql files. (%s)", folderName)
	}

	return nil
}
