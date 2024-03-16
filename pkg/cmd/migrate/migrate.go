package migrate

import (
	"fmt"
	"path"

	"github.com/marianop9/valkyrie-migrate/internal/constants"
	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
	"github.com/marianop9/valkyrie-migrate/pkg/valkyrie"
	"github.com/spf13/cobra"
)

func NewMigrateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate <migrationFolder> [dbName]",
		Short: "updates the database to the latest migration",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.MaximumNArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				migrationFolder, dbName string
			)

			if len(args) == 2 {
				migrationFolder, dbName = args[0], args[1]
			} else {
				migrationFolder, dbName = args[0], constants.DefaultDb
			}

			if path.Ext(dbName) != ".db" {
				return fmt.Errorf("invalid database file extension")
			}
			
			db, err := helpers.GetDb(dbName)
			if err != nil {
				return err
			}

			migrationRepo := repository.NewMigrationRepo(db)

			return valkyrie.NewMigrateApp(migrationRepo).Run(migrationFolder)
		},
	}

	return c
}
