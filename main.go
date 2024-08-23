package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/healthcheck"

	"github.com/liobrdev/simplepasswords_api_gateway/app"
	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/routes"
)

func main() {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		log.Fatalln("Failed to load config from environment:", err)
	}

	app := app.CreateApp(&conf)
	dbs := databases.Init(&conf)

	if err := dbs.ApiGateway.AutoMigrate(
		&models.User{},
		&models.ClientSession{},
		&models.MFAToken{},
	); err != nil {
		log.Fatalln("Failed api_gateway database auto-migrate:", err.Error())
	}

	if err := dbs.Logger.AutoMigrate(&models.Log{}); err != nil {
		log.Fatalln("Failed logger database auto-migrate:", err.Error())
	}

	app.Use(healthcheck.New())
	routes.Register(app, dbs, &conf)

	log.Fatal(app.Listen(conf.API_GATEWAY_HOST + ":" + conf.API_GATEWAY_PORT))
}
