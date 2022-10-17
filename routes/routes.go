package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/controllers/auth"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
)

func RegisterAPI(app *fiber.App, dbs *databases.Databases, conf *config.AppConfig) {
	api := app.Group("/api")
	auth.RegisterAuth(&api, dbs, conf)
}
