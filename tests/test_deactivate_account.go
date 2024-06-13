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

func testDeactivateAccount(t *testing.T, app *fiber.App, dbs *databases.Databases,
conf *config.AppConfig) {
	var clientIP string

	if conf.GO_FIBER_BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"email":"%s","password":"%s"}`

	setup.SetUpLogger(t, dbs)
	users, validTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

	validAuthHeader := "Token " + validTokens[0]

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "[]", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "[{}]", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(bodyFmt, "jane.doe@email.co", helpers.VALID_PW) + `]`

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "null", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "true", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "false", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, `"Valid JSON, but not an object."`, users[0].Slug, 400,
			utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, "{}", users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
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

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_email_400_bad_request", func(t *testing.T) {
		body := `{"email":null,"password":"` + helpers.VALID_PW + `"}`

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "", helpers.VALID_PW)

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
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

			testDeactivateAccountClientError(
				t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
				&models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.DeactivateAccount,
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

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_password_400_bad_request", func(t *testing.T) {
		body := `{"email":"jane.doe@email.co","password":null}`

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "jane.doe@email.co", "")

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_mismatch_token_and_slug_400_bad_request", func(t *testing.T) {
		email := users[0].EmailAddress
		body := fmt.Sprintf(bodyFmt, email, helpers.VALID_PW)

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[1].Slug, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.DeactivateAccount,
				Detail:          "urlParamSlug != user.Slug",
				Extra: 					 users[1].Slug + ":" + users[0].Slug,
				Level:           "error",
				Message:         utils.ErrorParams,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_failed_deactivate_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		
		email := users[0].EmailAddress
		body := fmt.Sprintf(bodyFmt, email, helpers.VALID_PW + "abc123")

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorFailedDeactivate, nil, nil,
			nil,
		)

		body = fmt.Sprintf(bodyFmt, "jake.doe@email.co", helpers.VALID_PW)

		testDeactivateAccountClientError(
			t, app, dbs, validAuthHeader, body, users[0].Slug, 400, utils.ErrorFailedDeactivate, nil, nil,
			nil,
		)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_200_ok", func(t *testing.T) {
		users, validTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		validAuthHeader := "Token " + validTokens[0]
		body := fmt.Sprintf(bodyFmt, users[0].EmailAddress, helpers.VALID_PW)

		testDeactivateAccountSuccess(t, app, dbs, validAuthHeader, body, users[0].Slug)
	})

	t.Run("valid_body_irrelevant_data_200_ok", func(t *testing.T) {
		users, validTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		validAuthHeader := "Token " + validTokens[0]
		body := fmt.Sprintf(
			`{"email":"%s","password":"%s","abc":123}`, users[0].EmailAddress, helpers.VALID_PW,
		)

		testDeactivateAccountSuccess(t, app, dbs, validAuthHeader, body, users[0].Slug)
	})
}

func testDeactivateAccountClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, body string,
	userSlug string, expectedStatus int, expectedDetail string,
	expectedFieldErrors map[string][]string, expectedNonFieldErrors []string, expectedLog *models.Log,
) {
	resp := newRequestDeactivateAccount(t, app, authHeader, body, userSlug)
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

func testDeactivateAccountSuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, authHeader string, body string,
	userSlug string,
) {
	var userCount int64
	helpers.CountUsers(t, dbs.ApiGateway, &userCount)
	require.EqualValues(t, 2, userCount)

	var sessionCount int64
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 8, sessionCount)

	resp := newRequestDeactivateAccount(t, app, authHeader, body, userSlug)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	helpers.CountUsers(t, dbs.ApiGateway, &userCount)
	require.EqualValues(t, 1, userCount)

	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 4, sessionCount)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		require.Empty(t, respBody)
	}
}

func newRequestDeactivateAccount(
	t *testing.T, app *fiber.App, authHeader string, body string, userSlug string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/auth/deactivate_account/" + userSlug, reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.DeactivateAccount)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
