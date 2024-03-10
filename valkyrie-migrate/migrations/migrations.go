package migrations

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/marianop9/valkyrie-migrate/valkyrie-migrate/repository"
)

const (
	dateFmt = "20060102"
)

func Migrate(repo *repository.MigrationRepo) {


}

func GetMigrationGroups(dirEntries []os.DirEntry) ([]*repository.MigrationGroup, error) {
	migrationGroups := make([]*repository.MigrationGroup, 0)

	for _, entry := range dirEntries {
		fmt.Printf("%+v\n", entry.Name())
		if entry.IsDir() {
			baseFolderName := entry.Name()
			entryParts := strings.Split(baseFolderName, "_")

			if len(entryParts) < 2 {
				return nil, errors.New("folder name doesn't match expected format")
			}

			date, err := time.Parse(dateFmt, entryParts[0]) // yyyymmdd
			if err != nil {
				return nil, errors.New("folder date doesn't match expected format")
			}

			group := repository.MigrationGroup{
				Name:       strings.Join(entryParts[1:], "_"),
				FolderName: baseFolderName,
				Date:       date,
			}

			migrationGroups = append(migrationGroups, &group)
		}
	}

	return migrationGroups, nil
}