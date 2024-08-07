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

func testAuthFirstFactor(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"email":"%s","password":"%s"}`

	setup.SetUpLogger(t, dbs)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testAuthFirstFactorClientError(
			t, app, dbs, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(bodyFmt, conf.ADMIN_EMAIL, helpers.VALID_PW) + `]`

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, "null", 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testAuthFirstFactorClientError(
			t, app, dbs, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, `"Valid JSON, but not an object."`, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testAuthFirstFactorClientError(
			t, app, dbs, "{}", 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("missing_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(`{"emial":"%s","password":"%s"}`, conf.ADMIN_EMAIL, helpers.VALID_PW)

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_email_400_bad_request", func(t *testing.T) {
		body := `{"email":null,"password":"` + helpers.VALID_PW + `"}`

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "", helpers.VALID_PW)

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("invalid_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "not-admin@test.com", helpers.VALID_PW)
		testAuthFirstFactorClientError(t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, nil)
	})

	t.Run("missing_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"email":"%s","passwrod":"%s"}`, conf.ADMIN_EMAIL, helpers.VALID_PW,
		)

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_password_400_bad_request", func(t *testing.T) {
		body := `{"email":"` + conf.ADMIN_EMAIL + `","password":null}`

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_EMAIL, "")

		testAuthFirstFactorClientError(
			t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthFirstFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_failed_login_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		setup.SetUpApiGatewayWithData(t, dbs, conf)

		body := fmt.Sprintf(bodyFmt, conf.ADMIN_EMAIL, helpers.VALID_PW + "abc123")
		testAuthFirstFactorClientError(t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, nil)

		body = fmt.Sprintf(bodyFmt, "not-admin@test.com", helpers.VALID_PW)
		testAuthFirstFactorClientError(t, app, dbs, body, 400, utils.ErrorFailedLogin, nil, nil, nil)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_200_ok", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_EMAIL, helpers.VALID_PW)
		testAuthFirstFactorSuccess(t, app, dbs, body, conf)
	})

	t.Run("valid_body_irrelevant_data_200_ok", func(t *testing.T) {
		validBodyIrrelevantData := fmt.Sprintf(
			`{"email":"%s","password":"%s","abc":123}`, conf.ADMIN_EMAIL, helpers.VALID_PW,
		)
		testAuthFirstFactorSuccess(t, app, dbs, validBodyIrrelevantData, conf)
	})
}

func testAuthFirstFactorClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestAuthFirstFactor(t, app, body)
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

func testAuthFirstFactorSuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, conf *config.AppConfig,
) {
	setup.SetUpApiGatewayWithData(t, dbs, conf)

	var mfaTokenCount int64
	helpers.CountMFATokens(t, dbs.ApiGateway, &mfaTokenCount)
	require.EqualValues(t, 4, mfaTokenCount)

	resp := newRequestAuthFirstFactor(t, app, body)
	require.Equal(t, 200, resp.StatusCode)

	helpers.CountMFATokens(t, dbs.ApiGateway, &mfaTokenCount)
	require.EqualValues(t, 3, mfaTokenCount)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var authFirstFactorRespBody controllers.AuthFirstFactorResponseBody

		if err := json.Unmarshal(respBody, &authFirstFactorRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Regexp(t, utils.TokenRegexp, authFirstFactorRespBody.MFAToken)

		var mfaToken models.MFAToken
		helpers.QueryTestMFATokenLatest(t, dbs.ApiGateway, &mfaToken)
		require.Equal(t, mfaToken.KeyDigest, utils.HashToken(authFirstFactorRespBody.MFAToken))
		require.Equal(t, mfaToken.OTPDigest, utils.HashToken(authFirstFactorRespBody.TestOTP))
	}
}

func newRequestAuthFirstFactor(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/auth/first_factor", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.AuthFirstFactor)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
