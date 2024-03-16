package migrate

import (
	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
	"github.com/marianop9/valkyrie-migrate/pkg/valkyrie"
	"github.com/spf13/cobra"
)

func NewMigrateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate",
		Short: "updates the database to the latest migration",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			migrationFolder, cnnString := args[0], args[1]
			
			db, err := helpers.GetDb(cnnString)
			if err != nil {
				return err
			}

			migrationRepo := repository.NewMigrationRepo(db)
			
			return valkyrie.NewMigrateApp(migrationRepo).Run(migrationFolder)
		},
	}

	return c;
}
