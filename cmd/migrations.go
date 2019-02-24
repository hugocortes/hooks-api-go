package cmd

import (
	"github.com/hugocortes/hooks-api/common/deps"
	"github.com/hugocortes/hooks-api/migrations"
	"github.com/spf13/cobra"
)

var migrationsCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Runs migrations and schema initialization",
	Long:  "This command runs all migrations required for hooks-api",
	Run: func(cmd *cobra.Command, args []string) {
		db := deps.Postgres()
		migrations.Run(db)
	},
}

func init() {
	RootCmd.AddCommand(migrationsCmd)

	deps.LoadEnv()
	deps.ConfigureLog()
}
