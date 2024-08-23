package databases

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
)

func Init(conf *config.AppConfig) *databases.Databases {
	var dbApiGateway *gorm.DB
	var dbLogger *gorm.DB

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	}

	dbApiGateway = openDbSession("./test_dbs/api_gateway.sqlite", &gormConfig)
	dbLogger = openDbSession("./test_dbs/logger.sqlite", &gormConfig)

	return &databases.Databases{ ApiGateway: dbApiGateway, Logger: dbLogger }
}

func openDbSession(dsn string, gormConfig *gorm.Config) (db *gorm.DB) {
	var err error

	if db, err = gorm.Open(sqlite.Open(dsn), gormConfig); err != nil {
		panic(err)
	}

	return
}
