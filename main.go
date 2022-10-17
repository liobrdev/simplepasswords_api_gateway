package main

import (
	"log"

	"github.com/liobrdev/simplepasswords_api_gateway/app"
	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func main() {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnvFile(&conf); err != nil {
		log.Fatalln("Failed to load config from '.env' file:", err)
	}

	app, dbs := app.CreateApp(&conf)

	if err := dbs.ApiGateway.AutoMigrate(
		&models.User{},
		&models.DeactivatedUser{},
		&models.ClientSession{},
	); err != nil {
		log.Fatalln("Failed api_gateway database auto-migrate:", err.Error())
	}

	if err := dbs.Logger.AutoMigrate(&models.Log{}); err != nil {
		log.Fatalln("Failed logger database auto-migrate:", err.Error())
	}

	app.Listen(conf.GO_FIBER_SERVER_HOST + conf.GO_FIBER_SERVER_PORT)
}
