package tests

import (
	"testing"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/app"
	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
)

func TestApp(t *testing.T) {
	var conf config.AppConfig

	if err := config.LoadConfigFromEnv(&conf); err != nil {
		t.Fatal("Failed to load config from environment:", err)
	}

	t.Run("is_behind_proxy", func(t *testing.T) {
		conf.GO_TESTING_CONTEXT = t
		conf.GO_FIBER_BEHIND_PROXY = true
		app, dbs := app.CreateApp(&conf)
		runTests(t, app, dbs, &conf)
	})

	t.Run("is_not_behind_proxy", func(t *testing.T) {
		conf.GO_TESTING_CONTEXT = t
		conf.GO_FIBER_BEHIND_PROXY = false
		app, dbs := app.CreateApp(&conf)
		runTests(t, app, dbs, &conf)
	})
}

func runTests(t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig) {
	t.Run("test_create_account", func(t *testing.T) {
		testCreateAccount(t, app, dbs, conf)
	})

	t.Run("test_log_in_account", func(t *testing.T) {
		testLogInAccount(t, app, dbs, conf)
	})

	t.Run("test_authorize_request", func(t *testing.T) {
		testAuthorizeRequest(t, app, dbs, conf)
	})

	t.Run("test_deactivate_account", func(t *testing.T) {
		testDeactivateAccount(t, app, dbs, conf)
	})

	t.Run("test_verify_email_try", func(t *testing.T) {
		testVerifyEmailTry(t, app, dbs, conf)
	})
}
