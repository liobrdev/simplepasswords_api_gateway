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

func testVerifyEmailTry(t *testing.T, app *fiber.App, dbs *databases.Databases,
conf *config.AppConfig) {
	var clientIP string

	if conf.GO_FIBER_BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	setup.SetUpLogger(t, dbs)
	users, validSessionTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

	authHeader := "Token " + validSessionTokens[0]

	t.Run("user_email_already_verified_204_no_content", func(t *testing.T) {
		users[0].EmailIsVerified = true
		dbs.ApiGateway.Save(&users[0])

		testVerifyEmailTryClientError(t, app, dbs, authHeader, 204, "", nil, nil, &models.Log{
			ClientIP:        clientIP,
			ClientOperation: utils.VerifyEmailTry,
			Detail:          "user.EmailIsVerified == true",
			Extra: 					 users[0].EmailAddress,
			Level:           "warn",
			Message:         utils.ErrorBadClient,
		})
	})

	t.Run("valid_try_204_no_content", func(t *testing.T) {
		users[0].EmailIsVerified = false
		dbs.ApiGateway.Save(&users[0])

		var emailTokenCount int64
		helpers.CountEmailVerificationTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 8, emailTokenCount)

		testVerifyEmailTrySuccess(t, app, authHeader)

		helpers.CountEmailVerificationTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 5, emailTokenCount)

		var emailToken models.EmailVerificationToken
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		tokenKey := emailToken.TokenKey

		testVerifyEmailTrySuccess(t, app, authHeader)

		helpers.CountEmailVerificationTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 5, emailTokenCount)

		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.Equal(t, tokenKey, emailToken.TokenKey)
	})
}

func testVerifyEmailTryClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestVerifyEmailTry(t, app, authHeader)
	require.Equal(t, expectedStatus, resp.StatusCode)

	if expectedStatus != 204 {
		helpers.AssertErrorResponseBody(t, resp, &utils.ErrorResponseBody{
			Detail:         expectedDetail,
			FieldErrors:    expectedFieldErrors,
			NonFieldErrors: expectedNonFieldErrors,
		})
	}

	if expectedLog != nil {
		var actualLog models.Log
		helpers.QueryTestLogLatest(t, dbs.Logger, &actualLog)
		helpers.AssertLog(t, expectedLog, &actualLog)
	}
}

func testVerifyEmailTrySuccess(t *testing.T, app *fiber.App, authHeader string) {
	resp := newRequestVerifyEmailTry(t, app, authHeader)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)	

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}
}

func newRequestVerifyEmailTry(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest("POST", "/api/auth/verify_email_try", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.VerifyEmailTry)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
