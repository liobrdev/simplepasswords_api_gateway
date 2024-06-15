package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func testVerifyEmailConfirm(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"token":"%s"}`

	setup.SetUpLogger(t, dbs)
	users, validSessionTokens, _, validEmailTokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

	authHeader := "Token " + validSessionTokens[0]

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(bodyFmt, "jane.doe@email.co", helpers.VALID_PW) + `]`

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "null", 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, `"Valid JSON, but not an object."`, 400, utils.ErrorBadRequest, nil,
			nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, "{}", 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("missing_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(`{"toekn":"%s"}`, validEmailTokens[1])

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_token_400_bad_request", func(t *testing.T) {
		body := `{"token":null}`

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "")

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("invalid_token_regexp_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, validEmailTokens[1][:79])

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Extra: 					 validEmailTokens[1][:79],
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)

		body = fmt.Sprintf(bodyFmt, validEmailTokens[1] + "a")

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Extra: 					 validEmailTokens[1] + "a",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)

		body = fmt.Sprintf(bodyFmt, validEmailTokens[1][:79] + "$")

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Extra: 					 validEmailTokens[1][:79] + "$",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)

		body = fmt.Sprintf(bodyFmt, validEmailTokens[1][:79] + " ")

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.VerifyEmailConfirm,
				Detail:          "token",
				Extra: 					 validEmailTokens[1][:79] + " ",
				Level:           "warn",
				Message:         utils.ErrorToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("no_match_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, validEmailTokens[3])

		testVerifyEmailConfirmClientError(
			t, app, dbs, authHeader, body, 400, utils.ErrorToken, nil, nil, nil,
		)
	})

	t.Run("user_email_already_verified_204_no_content", func(t *testing.T) {
		users[0].EmailIsVerified = true
		dbs.ApiGateway.Save(&users[0])
		body := fmt.Sprintf(bodyFmt, validEmailTokens[1])
		testVerifyEmailConfirmClientError(t, app, dbs, authHeader, body, 204, "", nil, nil, nil)
	})

	t.Run("valid_try_204_no_content", func(t *testing.T) {
		users[0].EmailIsVerified = false
		dbs.ApiGateway.Save(&users[0])

		var emailTokenCount int64
		helpers.CountEmailVerificationTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 8, emailTokenCount)

		body := fmt.Sprintf(bodyFmt, validEmailTokens[1])
		testVerifyEmailConfirmSuccess(t, app, authHeader, body)

		helpers.CountEmailVerificationTokens(t, dbs.ApiGateway, &emailTokenCount)
		require.EqualValues(t, 4, emailTokenCount)

		var user models.User
		helpers.QueryTestUserBySlug(t, dbs.ApiGateway, &user, users[0].Slug)
		require.Equal(t, true, user.EmailIsVerified)
	})
}

func testVerifyEmailConfirmClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, body string,
	expectedStatus int, expectedDetail string, expectedFieldErrors map[string][]string,
	expectedNonFieldErrors []string, expectedLog *models.Log,
) {
	resp := newRequestVerifyEmailConfirm(t, app, authHeader, body)
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

func testVerifyEmailConfirmSuccess(t *testing.T, app *fiber.App, authHeader string, body string) {
	resp := newRequestVerifyEmailConfirm(t, app, authHeader, body)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)	

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}
}

func newRequestVerifyEmailConfirm(
	t *testing.T, app *fiber.App, authHeader string, body string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/auth/verify_email_confirm", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.VerifyEmailConfirm)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
