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
	usersApi.Get("/", H.RetrieveUser)

	vaultsApi := api.Group("/vaults")
	vaultsApi.Post("/", H.VaultsCreateVault)
	vaultsApi.Get("/", H.VaultsListVaults)
	vaultsApi.Get("/:slug", H.VaultsRetrieveVault)
	vaultsApi.Patch("/:slug", H.VaultsUpdateVault)
	vaultsApi.Delete("/:slug", H.VaultsDeleteVault)

	entriesApi := api.Group("/entries")
	entriesApi.Post("/", H.VaultsCreateEntry)
	entriesApi.Get("/:slug", H.VaultsRetrieveEntry)
	entriesApi.Patch("/:slug", H.VaultsUpdateEntry)
	entriesApi.Delete("/:slug", H.VaultsDeleteEntry)

	secretsApi := api.Group("/secrets")
	secretsApi.Post("/", H.VaultsCreateSecret)
	secretsApi.Patch("/:slug", H.VaultsUpdateSecret)
	secretsApi.Delete("/:slug", H.VaultsDeleteSecret)
}
