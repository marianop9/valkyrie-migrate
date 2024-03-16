package cmd

import (
	initCmd "github.com/marianop9/valkyrie-migrate/pkg/cmd/init"
	"github.com/marianop9/valkyrie-migrate/pkg/cmd/migrate"
	"github.com/spf13/cobra"
)

func NewValkyrieCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "valkyrie",
		Short: "valkyrie-migrate is a tool for managing database migrations.",
	}

	rootCmd.AddCommand(
		migrate.NewMigrateCmd(),
		initCmd.NewInitCmd(),
	)

	return rootCmd
}
