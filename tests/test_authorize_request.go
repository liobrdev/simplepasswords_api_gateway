package tests

import (
	"io"
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

func testAuthorizeRequest(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.GO_FIBER_BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	setup.SetUpLogger(t, dbs)

	var dummyToken string
	var err error

	if dummyToken, err = utils.GenerateSlug(80); err != nil {
		t.Fatalf("Generate dummyToken failed: %s", err.Error())
	}

	t.Run("null_empty_token_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(t, app, dbs, "Token null", 401, utils.ErrorToken, nil, nil, nil)
		testAuthorizeRequestClientError(t, app, dbs, "token null", 401, utils.ErrorToken, nil, nil, nil)
		testAuthorizeRequestClientError(t, app, dbs, "Token ", 401, utils.ErrorToken, nil, nil, nil)
		testAuthorizeRequestClientError(t, app, dbs, "token ", 401, utils.ErrorToken, nil, nil, nil)
	})

	t.Run("invalid_token_regexp_401_unauthorized", func(t *testing.T) {
		testAuthorizeRequestClientError(
			t, app, dbs, "", 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "",
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "T0ken " + dummyToken, 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "T0ken " + dummyToken,
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "Bearer " + dummyToken, 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "Bearer " + dummyToken,
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, dummyToken, 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 dummyToken,
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + dummyToken + "abc", 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "Token " + dummyToken + "abc",
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + dummyToken[:79], 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "Token " + dummyToken[:79],
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token a " + dummyToken[2:80], 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "Token a " + dummyToken[2:80],
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + dummyToken[:79] + "!", 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "Authorization",
				Extra:					 "Token " + dummyToken[:79] + "!",
				Level:           "error",
				Message:         utils.ErrorToken,
			},
		)
	})

	t.Run("valid_token_no_match_401_unauthorized", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + dummyToken, 401, utils.ErrorToken, nil, nil, nil,
		)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_token_ip_mismatch_401_unauthorized", func(t *testing.T) {
		_, validTokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + validTokens[1], 401, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.TestAuthReq,
				Detail:          "c.IP() != thisSession.ClientIP",
				Extra:					 clientIP + " != " + helpers.OLD_IP,
				Level:           "error",
				Message:         utils.ErrorIPMismatch,
			},
		)
	})

	t.Run("valid_token_expired_401_unauthorized", func(t *testing.T) {
		_, _, expiredTokens := setup.SetUpApiGatewayWithData(t, dbs, conf)
		setup.SetUpLogger(t, dbs)

		var sessionCount int64
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 8, sessionCount)

		testAuthorizeRequestClientError(
			t, app, dbs, "Token " + expiredTokens[0], 401, utils.ErrorToken, nil, nil, nil,
		)

		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 6, sessionCount)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_slug_204_no_content", func(t *testing.T) {
		_, validTokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		var sessionCount int64
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 8, sessionCount)

		testAuthorizeRequestSuccess(t, app, "Token " + validTokens[0])
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 6, sessionCount)

		testAuthorizeRequestSuccess(t, app, "Token " + validTokens[2])
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 4, sessionCount)
	})
}

func testAuthorizeRequestClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestAuthorizeRequest(t, app, authHeader)
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

func testAuthorizeRequestSuccess(t *testing.T, app *fiber.App, authHeader string) {
	resp := newRequestAuthorizeRequest(t, app, authHeader)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}
}

func newRequestAuthorizeRequest(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest("GET", "/api/auth/restricted", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.TestAuthReq)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
