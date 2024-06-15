package databases

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
)

type Databases struct {
	ApiGateway *gorm.DB
	Logger     *gorm.DB
}

func Init(conf *config.AppConfig) *Databases {
	var dbApiGateway *gorm.DB
	var dbLogger *gorm.DB

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	if conf.ENVIRONMENT != "production" {
		dbApiGateway = openDbSession("sqlite", "./test_dbs/api_gateway.sqlite", &gormConfig)
		dbLogger = openDbSession("sqlite", "./test_dbs/logger.sqlite", &gormConfig)
	} else {
		dbApiGateway = openDbSession("postgres", fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
			conf.API_GATEWAY_DB_USER,
			conf.API_GATEWAY_DB_PASSWORD,
			conf.API_GATEWAY_DB_HOST,
			conf.API_GATEWAY_DB_PORT,
			conf.API_GATEWAY_DB_NAME,
		), &gormConfig)

		dbLogger = openDbSession("postgres", fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
			conf.LOGGER_DB_USER,
			conf.LOGGER_DB_PASSWORD,
			conf.LOGGER_DB_HOST,
			conf.LOGGER_DB_PORT,
			conf.LOGGER_DB_NAME,
		), &gormConfig)
	}

	return &Databases{dbApiGateway, dbLogger}
}

func openDbSession(driver string, dsn string, gormConfig *gorm.Config) (db *gorm.DB) {
	var err error

	if driver == "sqlite" {
		if db, err = gorm.Open(sqlite.Open(dsn), gormConfig); err != nil {
			panic(err)
		}
	} else if driver == "postgres" {
		if db, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
			panic(err)
		}
	} else {
		panic("Unsupported database driver: " + driver)
	}

	return
}
