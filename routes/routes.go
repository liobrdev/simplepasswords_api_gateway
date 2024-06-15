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
	authApi.Post("/log_in_account", H.LogInAccount)

	if H.Conf.ENVIRONMENT == "testing" {
		authApi.Get("/restricted", H.AuthorizeRequest, H.Restricted)
	}

	authApi.Post("/deactivate_account/:slug", H.AuthorizeRequest, H.DeactivateAccount)
	authApi.Post("/verify_email_try", H.AuthorizeRequest, H.VerifyEmailTry)
	authApi.Post("/verify_email_confirm", H.AuthorizeRequest, H.VerifyEmailConfirm)
}
