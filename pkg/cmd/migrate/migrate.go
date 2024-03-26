package migrate

import (
	"fmt"
	"path"
	"strings"

	"github.com/marianop9/valkyrie-migrate/internal/constants"
	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/models"
	postgresRepo "github.com/marianop9/valkyrie-migrate/internal/repository/postgres"
	sqliteRepo "github.com/marianop9/valkyrie-migrate/internal/repository/sqlite"
	"github.com/marianop9/valkyrie-migrate/pkg/valkyrie"
	"github.com/spf13/cobra"
)

func NewMigrateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate <migrationFolder> [dbName]",
		Short: "updates the database to the latest migration",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				migrationFolder, dbName string
			)

			if len(args) == 2 {
				migrationFolder, dbName = args[0], args[1]
			} else {
				migrationFolder, dbName = args[0], constants.DefaultDb
			}

			var migrationRepo models.MigrationStorer

			if strings.HasPrefix(dbName, "postgresql://") {
				db, err := helpers.GetPostgresDb(dbName)
				if err != nil {
					return err
				}
				migrationRepo = postgresRepo.NewMigrationRepo(db)

			} else if path.Ext(dbName) == ".db" {
				db, err := helpers.GetDb(dbName)
				if err != nil {
					return err
				}
				migrationRepo = sqliteRepo.NewMigrationRepo(db)
			} else {
				return fmt.Errorf("invalid database file extension")
			}

			return valkyrie.NewMigrateApp(migrationRepo).Run(migrationFolder)
		},
	}

	return c
}
