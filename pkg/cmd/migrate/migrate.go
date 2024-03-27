package migrate

import (
	"errors"
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

var ErrNoMigrationFolder = errors.New("the folder containing migrations must be specified")

func NewMigrateCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "migrate <migrationFolder> [connFile]",
		Short: "Updates the database to the latest migration",
		Long:  "Updates the database to the latest migration. To specify a database, pass the path to the connFile as the second argument, or specify the connection directly with --conn",
		Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(2)),
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				migrationFolder, connString string
			)

			connString, err := cmd.Flags().GetString(constants.ConnFlagName)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return ErrNoMigrationFolder
			} else {
				migrationFolder = args[0]
			}

			if connString == "" {
				if len(args) > 1 {
					if connString, err = helpers.GetConnString(args[1]); err != nil {
						return err
					}
				} else {
					connString = constants.DefaultDb
				}
			}

			var migrationRepo models.MigrationStorer

			if strings.HasPrefix(connString, "postgresql://") {
				db, err := helpers.GetPostgresDb(connString)
				if err != nil {
					return err
				}
				migrationRepo = postgresRepo.NewMigrationRepo(db)

			} else if path.Ext(connString) == ".db" {
				db, err := helpers.GetDb(connString)
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

	c.PersistentFlags().String(constants.ConnFlagName, "", "directly specifies a db connection, ignoring the config file")

	return c
}
