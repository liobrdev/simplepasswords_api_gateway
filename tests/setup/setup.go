package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func SetUpApiGateway(t *testing.T, dbs *databases.Databases) {
	TearDownApiGateway(t, dbs)

	if err := dbs.ApiGateway.AutoMigrate(&models.User{}, &models.ClientSession{}); err != nil {
		t.Fatalf("Failed database auto-migrate: %s", err.Error())
	}
}

func SetUpApiGatewayWithData(t *testing.T, dbs *databases.Databases, conf *config.AppConfig) (
	users []models.User,
	validTokens []string,
	expiredTokens []string,
) {
	SetUpApiGateway(t, dbs)
	users, validTokens, expiredTokens = populateTestDBApiGateway(t, dbs, conf)
	return
}

func SetUpLogger(t *testing.T, dbs *databases.Databases) {
	TearDownLogger(t, dbs)

	if err := dbs.Logger.AutoMigrate(&models.Log{}); err != nil {
		t.Fatalf("Failed database auto-migrate: %s", err.Error())
	}
}

func TearDownApiGateway(t *testing.T, dbs *databases.Databases) {
	if result := dbs.ApiGateway.Exec(
		"DROP TABLE IF EXISTS users",
	); result.Error != nil {
		t.Fatalf("Test database tear-down failed: %s", result.Error.Error())
	}

	if result := dbs.ApiGateway.Exec(
		"DROP TABLE IF EXISTS client_sessions",
	); result.Error != nil {
		t.Fatalf("Test database tear-down failed: %s", result.Error.Error())
	}
}

func TearDownLogger(t *testing.T, dbs *databases.Databases) {
	if result := dbs.Logger.Exec("DROP TABLE IF EXISTS logs"); result.Error != nil {
		t.Fatalf("Test database tear-down failed: %s", result.Error.Error())
	}
}
