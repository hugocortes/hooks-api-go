package cmd

import (
	"os"

	_binsHandlers "github.com/hugocortes/hooks-api/bins/handlers"
	_binsInterfaces "github.com/hugocortes/hooks-api/bins/interfaces"
	_binsRepository "github.com/hugocortes/hooks-api/bins/repository"
	"github.com/hugocortes/hooks-api/common/deps"
	"github.com/hugocortes/hooks-api/common/middleware"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the server",
	Long:  "This command starts the server",
	Run: func(cmd *cobra.Command, args []string) {
		deps.LoadEnv()
		deps.ConfigureLog()

		router := deps.Router()
		postgres := deps.Postgres()
		redis := deps.Redis()

		// Bin initialization
		binRepo := _binsRepository.New(postgres, redis)
		binHandler := _binsHandlers.New(binRepo)
		binInter := _binsInterfaces.New(binHandler)
		binInter.AddRoutes(router)

		// start http
		middle := middleware.New(router)
		router.NoRoute(middle.NotFound)
		router.Use(middle.CorsConfig())
		middle.Auth()
		router.Run(":" + os.Getenv("PORT"))
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
}
