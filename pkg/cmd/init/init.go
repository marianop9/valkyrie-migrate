package init

import (
	"fmt"
	"path"

	"github.com/marianop9/valkyrie-migrate/internal/constants"
	"github.com/marianop9/valkyrie-migrate/pkg/valkyrie"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Creates or verifies the connectino to the database.",
		Long:  "Creates or pings the specified database. The deafult database name is used if none is specificed.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var dbName string

			if len(args) > 0 {
				dbName = args[0]
			} else {
				dbName = constants.DefaultDb
			}

			if path.Ext(dbName) != ".db" {
				return fmt.Errorf("invalid database file extension")
			}

			return valkyrie.Init(dbName)
		},
	}

	return c
}
