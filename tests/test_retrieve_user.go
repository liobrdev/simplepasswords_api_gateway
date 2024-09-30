package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/setup"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func testRetrieveUser(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	t.Run("wrong_client_operation_header_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)

		testRetrieveUserClientError(
			t, app, dbs, "wrong_operation", "Token " + validTokens[0], 400, utils.ErrorBadRequest,
			nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.RetrieveUser,
				Detail:          "wrong_operation",
				Level:           "warn",
				Message:         utils.ErrorClientOperation,
			},
		)
	})

	t.Run("valid_token_200_ok", func(t *testing.T) {
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)

		testRetrieveUserSuccess(
			t, app, utils.RetrieveUser, user.Slug, user.Name, "Token " + validTokens[0],
		)
	})
}

func testRetrieveUserClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, clientOperation, authHeader string,
	expectedStatus int, expectedDetail string, expectedFieldErrors map[string][]string,
	expectedNonFieldErrors []string, expectedLog *models.Log,
) {
	resp := newRequestRetrieveUser(t, app, clientOperation, authHeader)
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

func testRetrieveUserSuccess(
	t *testing.T, app *fiber.App, clientOperation, slug, name, authHeader string,
) {
	resp := newRequestRetrieveUser(t, app, clientOperation, authHeader)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var user models.User

		if err := json.Unmarshal(respBody, &user); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, slug, user.Slug)
		require.Equal(t, name, user.Name)
		require.True(t, user.EmailIsVerified)
		require.True(t, user.PhoneIsVerified)
	}
}

func newRequestRetrieveUser(
	t *testing.T, app *fiber.App, clientOperation, authHeader string,
) *http.Response {
	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", clientOperation)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
