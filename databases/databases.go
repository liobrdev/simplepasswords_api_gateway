package databases

import (
	"fmt"

	"gorm.io/driver/postgres"
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

	dbApiGateway = openDbSession(fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
		conf.API_GATEWAY_DB_USER,
		conf.API_GATEWAY_DB_PASSWORD,
		conf.API_GATEWAY_DB_HOST,
		conf.API_GATEWAY_DB_PORT,
		conf.API_GATEWAY_DB_NAME,
	), &gormConfig)

	dbLogger = openDbSession(fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=UTC",
		conf.LOGGER_DB_USER,
		conf.LOGGER_DB_PASSWORD,
		conf.LOGGER_DB_HOST,
		conf.LOGGER_DB_PORT,
		conf.LOGGER_DB_NAME,
	), &gormConfig)

	return &Databases{dbApiGateway, dbLogger}
}

func openDbSession(dsn string, gormConfig *gorm.Config) (db *gorm.DB) {
	var err error
	
	if db, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
		panic(err)
	}

	return
}
