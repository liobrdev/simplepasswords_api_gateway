package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/controllers"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/setup"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func testVerifyEmailTry(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	t.Run("user_email_already_verified_200_ok", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)

		testVerifyEmailTryClientError(
			t, app, dbs, "Token " + validTokens[0], 200, "", nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailTry,
				Level:           "warn",
				Message:         utils.ErrorAlreadyVerified,
				UserSlug:				 user.Slug,
			},
		)
	})

	t.Run("retry_too_soon_valid_tokens_200_ok", func(t *testing.T) {
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)
		user.EmailIsVerified = false
		dbs.ApiGateway.Save(&user)

		tokenString, _ := utils.GenerateSlug(80)
		oneTimePasscode, _ := utils.GenerateOTP()
		now := time.Now().UTC()

		dbs.ApiGateway.Create(&models.EmailVerificationToken{
			UserSlug:  user.Slug,
			KeyDigest: utils.HashToken(tokenString),
			OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
			TokenKey:	 tokenString[:16],
			CreatedAt: now,
			ExpiresAt: now.Add(time.Duration(10) * time.Minute),
		})

		var emailTokenCount int64
		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 1, emailTokenCount)

		testVerifyEmailTryClientError(
			t, app, dbs, "Token " + validTokens[0], 200, "", nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailTry,
				Level:           "warn",
				Message:         "Too soon retry",
				UserSlug:				 user.Slug,
			},
		)

		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 1, emailTokenCount)

		var emailToken models.EmailVerificationToken
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.Equal(t, tokenString[:16], emailToken.TokenKey)
	})

	t.Run("cleanup_expired_tokens_200_ok", func(t *testing.T) {
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)
		user.EmailIsVerified = false
		dbs.ApiGateway.Save(&user)
		expiredEmailTokens := setup.CreateExpiredTestEmailVerificationTokens(&user, t, dbs)

		var emailToken models.EmailVerificationToken
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.Equal(t, expiredEmailTokens[1].TestEmailToken[:16], emailToken.TokenKey)

		var emailTokenCount int64
		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 2, emailTokenCount)

		testVerifyEmailTrySuccess(t, app, dbs, "Token " + validTokens[0], true)
		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 1, emailTokenCount)

		emailToken = models.EmailVerificationToken{}
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.NotEqual(t, expiredEmailTokens[1].TestEmailToken[:16], emailToken.TokenKey)
	})

	t.Run("cleanup_older_valid_tokens_200_ok", func(t *testing.T) {
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)
		user.EmailIsVerified = false
		dbs.ApiGateway.Save(&user)
		setup.CreateExpiredTestEmailVerificationTokens(&user, t, dbs)
		validEmailTokens := setup.CreateValidTestEmailVerificationTokens(&user, t, dbs)

		var emailToken models.EmailVerificationToken
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.Equal(t, validEmailTokens[1].TestEmailToken[:16], emailToken.TokenKey)

		var emailTokenCount int64
		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 4, emailTokenCount)

		testVerifyEmailTrySuccess(t, app, dbs, "Token " + validTokens[0], true)
		helpers.CountEmailTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 1, emailTokenCount)

		emailToken = models.EmailVerificationToken{}
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.NotEqual(t, validEmailTokens[1].TestEmailToken[:16], emailToken.TokenKey)
	})
}

func testVerifyEmailTryClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestVerifyEmailTry(t, app, authHeader)
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

func testVerifyEmailTrySuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, bodyHasTestData bool,
) {
	resp := newRequestVerifyEmailTry(t, app, authHeader)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else if bodyHasTestData {
		var verifyEmailTryRespBody controllers.VerifyEmailTryResponseBody

		if err := json.Unmarshal(respBody, &verifyEmailTryRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		var emailToken models.EmailVerificationToken
		helpers.QueryTestEmailTokenLatest(t, dbs.ApiGateway, &emailToken)
		require.Equal(t, emailToken.TokenKey, verifyEmailTryRespBody.TestEmailToken[:16])
		require.Equal(t, emailToken.OTPDigest, utils.HashToken(verifyEmailTryRespBody.TestOTP))
	} else {		
		require.Equal(t, "{}", string(respBody))
	}
}

func newRequestVerifyEmailTry(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest(http.MethodPost, "/api/auth/verify_email_try", nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Client-Operation", utils.VerifyEmailTry)
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
