package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/app"
	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/routes"
	testDBs "github.com/liobrdev/simplepasswords_api_gateway/tests/databases"
)

func TestApp(t *testing.T) {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		t.Fatal("Failed to load config from environment:", err)
	}

	t.Run("is_behind_proxy", func(t *testing.T) {
		conf.ENVIRONMENT = "testing"
		conf.GO_TESTING_CONTEXT = t
		conf.BEHIND_PROXY = true
		app := app.CreateApp(&conf)
		dbs := testDBs.Init(&conf)
		routes.Register(app, dbs, &conf)
		runTests(t, app, dbs, &conf)
	})

	t.Run("is_not_behind_proxy", func(t *testing.T) {
		conf.ENVIRONMENT = "testing"
		conf.GO_TESTING_CONTEXT = t
		conf.BEHIND_PROXY = false
		app := app.CreateApp(&conf)
		dbs := testDBs.Init(&conf)
		routes.Register(app, dbs, &conf)
		runTests(t, app, dbs, &conf)
	})
}

func runTests(t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig) {
	t.Run("test_create_account", func(t *testing.T) {
		testCreateAccount(t, app, dbs, conf)
	})

	t.Run("test_auth_first_factor", func(t *testing.T) {
		testAuthFirstFactor(t, app, dbs, conf)
	})
	
	t.Run("test_auth_second_factor", func(t *testing.T) {
		testAuthSecondFactor(t, app, dbs, conf)
	})

	t.Run("test_authorize_request", func(t *testing.T) {
		testAuthorizeRequest(t, app, dbs, conf)
	})

	t.Run("test_retrieve_user", func(t *testing.T) {
		testRetrieveUser(t, app, dbs, conf)
	})

	t.Run("test_client_operation_header", func(t *testing.T) {
		testClientOperationHeader(t, app, dbs, conf)
	})
}
