package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/controllers"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
)

func Register(app *fiber.App, dbs *databases.Databases, conf *config.AppConfig) {
	H := controllers.Handler{DBs: dbs, Conf: conf}
	api := app.Group("/api")

	authApi := api.Group("/auth")
	authApi.Post("/create_account", H.CreateAccount)
	authApi.Post("/first_factor", H.AuthFirstFactor)
	authApi.Post("/second_factor", H.AuthSecondFactor)

	app.Use(H.AuthorizeRequest)

	if H.Conf.ENVIRONMENT == "testing" {
		authApi.Get("/restricted", H.Restricted)
	}

	usersApi := api.Group("/users")
	usersApi.Get("/", H.VaultsRetrieveUser)

	// vaultsApi := api.Group("/vaults")
	// vaultsApi.Post("/", H.CreateVault)
	// vaultsApi.Get("/:slug", H.RetrieveVault)
	// vaultsApi.Patch("/:slug", H.UpdateVault)
	// vaultsApi.Delete("/:slug", H.DeleteVault)

	// entriesApi := api.Group("/entries")
	// entriesApi.Post("/", H.CreateEntry)
	// entriesApi.Get("/:slug", H.RetrieveEntry)
	// entriesApi.Patch("/:slug", H.UpdateEntry)
	// entriesApi.Delete("/:slug", H.DeleteEntry)

	// secretsApi := api.Group("/secrets")
	// secretsApi.Post("/", H.CreateSecret)
	// secretsApi.Patch("/:slug", H.UpdateSecret)
	// secretsApi.Delete("/:slug", H.DeleteSecret)
}
