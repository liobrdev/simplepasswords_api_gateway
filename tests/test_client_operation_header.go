package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/setup"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func testClientOperationHeader(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	dummySlug := helpers.NewSlug(t)
	
	clientOperations := map[string][]string{
		"create_account":			{"POST", "/api/auth/create_account"},
		"auth_first_factor":	{"POST", "/api/auth/first_factor"},
		"auth_second_factor":	{"POST", "/api/auth/second_factor"},
		"logout_account":			{"POST", "/api/auth/logout_account"},
		"retrieve_user":			{"GET", "/api/users"},
		"create_vault":				{"POST", "/api/vaults"},
		"list_vaults":				{"GET", "/api/vaults"},
		"retrieve_vault":			{"GET", "/api/vaults/" + dummySlug},
		"update_vault":				{"PATCH", "/api/vaults/" + dummySlug},
		"delete_vault":				{"DELETE", "/api/vaults/" + dummySlug},
		"create_entry":				{"POST", "/api/entries"},
		"retrieve_entry":			{"GET", "/api/entries/" + dummySlug},
		"update_entry":				{"PATCH", "/api/entries/" + dummySlug},
		"delete_entry":				{"DELETE", "/api/entries/" + dummySlug},
		"create_secret":			{"POST", "/api/secrets"},
		"update_secret":			{"PATCH", "/api/secrets/" + dummySlug},
		"delete_secret":			{"DELETE", "/api/secrets/" + dummySlug},
	}

	setup.SetUpLogger(t, dbs)
	_, validSessionTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

	for operation, route := range clientOperations {
		t.Run("wrong_client_operation_" + operation + "_400_bad_request", func(t *testing.T) {
			testClientOperationHeaderError(
				t, app, dbs, conf, route[0], route[1], "Token " + validSessionTokens[0], "wrong_operation",
				400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: operation,
					Detail:          "wrong_operation",
					Level:           "warn",
					Message:         utils.ErrorClientOperation,
				},
			)
		})
	}
}

func testClientOperationHeaderError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
	method, target, authHeader, clientOperation string,
	expectedStatus int, expectedDetail string, expectedFieldErrors map[string][]string,
	expectedNonFieldErrors []string, expectedLog *models.Log,
) {
	resp := newRequestGeneric(t, app, conf, method, target, authHeader, clientOperation)
	require.Equal(t, expectedStatus, resp.StatusCode)

	helpers.AssertErrorResponseBody(t, resp, &utils.ErrorResponseBody{
		Detail:         expectedDetail,
		FieldErrors:    expectedFieldErrors,
		NonFieldErrors: expectedNonFieldErrors,
	})

	if expectedLog != nil {
		var actualLog models.Log
		helpers.QueryTestLogLatest(t, dbs.Logger, &actualLog)
		helpers.AssertLog(t, expectedLog, &actualLog)
	}
}

func newRequestGeneric(
	t *testing.T, app *fiber.App, conf *config.AppConfig,
	method, target, authHeader, clientOperation string,
) *http.Response {
	req := httptest.NewRequest(method, target, nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Client-Operation", clientOperation)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set(conf.PASSWORD_HEADER_KEY, helpers.HexHash1)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
