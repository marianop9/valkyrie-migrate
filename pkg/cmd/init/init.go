package init

import (
	"encoding/json"
	"os"

	"github.com/marianop9/valkyrie-migrate/internal/constants"
	"github.com/marianop9/valkyrie-migrate/pkg/valkyrie"
	"github.com/spf13/cobra"
)

var connFlag = "conn"

func NewInitCmd() *cobra.Command {
	c := &cobra.Command{
		Use:       "init config-file [--conn db-connection]",
		Short:     "Creates or verifies the connectino to the database.",
		Long:      "Creates or pings the specified database. The deafult database is used if no config file is specificed.",
		Args:      cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			connString, err := cmd.Flags().GetString(connFlag)
			if err != nil  {
				return err
			} 
			
			if connString != "" {
				return valkyrie.Init(connString)
			}

			if len(args) > 0 {
				connFilePath := args[0]

				if connString, err = getConnString(connFilePath); err != nil {
					return err
				}
			} else {
				connString = constants.DefaultDb
			}

			return valkyrie.Init(connString)
		},
	}

	c.PersistentFlags().String(connFlag, "", "directly specifies a db connection, ignoring the config file")

	return c
}

type ConnFile struct {
	ConnectionString string
}

func getConnString(connFilePath string) (string, error) {
	buf, err := os.ReadFile(connFilePath)

	connFile := ConnFile{}

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(buf, &connFile)

	return connFile.ConnectionString, err
}