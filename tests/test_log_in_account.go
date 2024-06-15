package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func testLogInAccount(t *testing.T, app *fiber.App, dbs *databases.Databases,
conf *config.AppConfig) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"email":"%s","password":"%s"}`

	setup.SetUpLogger(t, dbs)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testLogInAccountClientError(
			t, app, dbs, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(bodyFmt, "jane.doe@email.co", helpers.VALID_PW) + `]`

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, "null", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testLogInAccountClientError(
			t, app, dbs, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, `"Valid JSON, but not an object."`, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testLogInAccountClientError(
			t, app, dbs, "{}", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("missing_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"emial":"%s","password":"%s"}`, "spelled@wrong.co", helpers.VALID_PW,
		)

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_email_400_bad_request", func(t *testing.T) {
		body := `{"email":null,"password":"` + helpers.VALID_PW + `"}`

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "", helpers.VALID_PW)

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_long_email_400_bad_request", func(t *testing.T) {
		if email, err := utils.GenerateSlug(257); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			body := fmt.Sprintf(bodyFmt, email, helpers.VALID_PW)

			testLogInAccountClientError(
				t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.LogInAccount,
					Detail:          email,
					Level:           "warn",
					Message:         utils.ErrorAcctEmail,
					RequestBody:     body,
				},
			)
		}
	})

	t.Run("missing_password_400_bad_request", func(t *testing.T) {
		body := `{"email":"jane.doe@email.co","passwrod":"$pelledWr0ng1234"}`

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_password_400_bad_request", func(t *testing.T) {
		body := `{"email":"jane.doe@email.co","password":null}`

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "jane.doe@email.co", "")

		testLogInAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.LogInAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_failed_login_404_not_found", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		users, _, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		email := users[0].EmailAddress
		body := fmt.Sprintf(bodyFmt, email, helpers.VALID_PW + "abc123")

		testLogInAccountClientError(t, app, dbs, body, 404, utils.ErrorFailedLogin, nil, nil, nil)

		body = fmt.Sprintf(bodyFmt, "jake.doe@email.co", helpers.VALID_PW)

		testLogInAccountClientError(t, app, dbs, body, 404, utils.ErrorFailedLogin, nil, nil, nil)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_200_ok", func(t *testing.T) {
		users, _, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		email := users[0].EmailAddress
		body := fmt.Sprintf(bodyFmt, email, helpers.VALID_PW)
		testLogInAccountSuccess(t, app, dbs, body, email)
	})

	t.Run("valid_body_irrelevant_data_200_ok", func(t *testing.T) {
		users, _, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		email := users[0].EmailAddress
		validBodyIrrelevantData := fmt.Sprintf(
			`{"email":"%s","password":"%s","abc":123}`, email, helpers.VALID_PW,
		)
		testLogInAccountSuccess(t, app, dbs, validBodyIrrelevantData, email)
	})
}

func testLogInAccountClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestLogInAccount(t, app, body)
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

func testLogInAccountSuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, email string,
) {
	var sessionCount int64
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 8, sessionCount)

	resp := newRequestLogInAccount(t, app, body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 9, sessionCount)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var logInAcctRespBody controllers.LogInAccountResponseBody

		if err := json.Unmarshal(respBody, &logInAcctRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Regexp(t, utils.TokenRegexp, logInAcctRespBody.Token)

		var user models.User
		helpers.QueryTestUserByEmail(t, dbs.ApiGateway, &user, email)
		require.Equal(t, user.Slug, logInAcctRespBody.User.Slug)
		require.Equal(t, user.Name, logInAcctRespBody.User.Name)
		require.Equal(t, user.EmailAddress, logInAcctRespBody.User.EmailAddress)
		require.Equal(t, user.EmailIsVerified, logInAcctRespBody.User.EmailIsVerified)
		require.Equal(t, false, logInAcctRespBody.User.EmailIsVerified)
		require.Equal(t, user.PhoneNumber, logInAcctRespBody.User.PhoneNumber)
		require.Equal(t, "", logInAcctRespBody.User.PhoneNumber)
		require.Equal(t, user.PhoneIsVerified, logInAcctRespBody.User.PhoneIsVerified)
		require.Equal(t, false, logInAcctRespBody.User.PhoneIsVerified)
		require.Equal(t, user.MfaIsEnabled, logInAcctRespBody.User.MfaIsEnabled)
		require.Equal(t, false, logInAcctRespBody.User.MfaIsEnabled)
		require.Equal(t, user.IsActive, logInAcctRespBody.User.IsActive)
		require.Equal(t, true, logInAcctRespBody.User.IsActive)
		require.Empty(t, logInAcctRespBody.User.CreatedAt)
		require.Empty(t, logInAcctRespBody.User.UpdatedAt)
		require.Empty(t, logInAcctRespBody.User.PasswordSalt)
		require.Empty(t, logInAcctRespBody.User.PasswordHash)

		var session models.ClientSession
		helpers.QueryTestClientSessionLatest(t, dbs.ApiGateway, &session)
		require.Equal(t, session.UserSlug, logInAcctRespBody.User.Slug)
		require.Equal(t, session.TokenKey, logInAcctRespBody.Token[:16])
		require.Equal(t, session.Digest, utils.HashToken(logInAcctRespBody.Token))
	}
}

func newRequestLogInAccount(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/auth/log_in_account", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.LogInAccount)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
