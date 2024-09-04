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

func testCreateAccount(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"name":"%s","email":"%s","phone":"%s","password":"%s","password_2":"%s"}`

	setup.SetUpLogger(t, dbs)

	t.Run("wrong_client_operation_header_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, "wrong_operation", body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "wrong_operation",
				Level:           "warn",
				Message:         utils.ErrorClientOperation,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(
			bodyFmt,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		) + `]`

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "null", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, `"Valid JSON, but not an object."`, 400,
			utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, "{}", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("invalid_name_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		body := fmt.Sprintf(
			bodyFmt, "NotAdmin", conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, nil,
		)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("missing_name_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"nmae":"%s","email":"%s","phone":"%s","password":"%s","password_2":"%s"}`,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_name_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":null,"email":"%s","phone":"%s","password":"%s","password_2":"%s"}`,
			conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_name_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, "", conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","emial":"%s","phone":"%s","password":"%s","password_2":"%s"}`,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":null,"phone":"%s","password":"%s","password_2":"%s"}`,
			conf.ADMIN_NAME, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, conf.ADMIN_NAME, "", conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("invalid_email_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)

		body := fmt.Sprintf(
			bodyFmt,
			conf.ADMIN_NAME, "wrong@test.com", conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, nil,
		)
	})

	t.Run("missing_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":"%s","passwrod":"%s","password_2":"%s"}`,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("unmatching_passwords_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW,
			"N0nmatchpa$sw0rd",
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorNonMatchPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("password_is_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, conf.ADMIN_EMAIL,
			conf.ADMIN_EMAIL,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Password is email",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_long_password_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		pwLength := 72 - len(conf.ADMIN_SALT_1) - len(conf.ADMIN_SALT_2)

		if slug, err := utils.GenerateSlug(pwLength - 3); err != nil {
			t.Fatalf("Generate long password failed: %s", err.Error())
		} else {
			pw := slug + "qQ1!"
			body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)
			testCreateAccountClientError(
				t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, nil,
			)
		}
	
		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("too_short_password_400_bad_request", func(t *testing.T) {
		if slug, err := utils.GenerateSlug(11); err != nil {
			t.Fatalf("Generate short password failed: %s", err.Error())
		} else {
			pw := slug + "qQ1!"
			body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)
			testCreateAccountClientError(
				t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, nil,
			)
		}

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("password_missing_uppercase_400_bad_request", func(t *testing.T) {
		pw := "n0uppercasepa$sw0rd"
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Missing uppercase",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("password_missing_lowercase_400_bad_request", func(t *testing.T) {
		pw := "N0LOWERCASEPA$SW0RD"
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Missing lowercase",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("password_missing_number_400_bad_request", func(t *testing.T) {
		pw := "nONumberPA$SWoRD"
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Missing number",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("password_missing_special_char_400_bad_request", func(t *testing.T) {
		pw := "nOSp3cialCharPA5SWoRD"
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Missing special char",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("password_has_whitespace_400_bad_request", func(t *testing.T) {
		pw := "h@s Spa(3	PA5SWoRD"
		body := fmt.Sprintf(bodyFmt, conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, pw, pw)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Has whitespace",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_user_already_exists_409_conflict", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		setup.SetUpApiGatewayWithData(t, dbs, conf)

		body := fmt.Sprintf(
			bodyFmt,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 500, utils.ErrorServer, nil, nil, nil,
		)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_201_created", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountSuccess(t, app, dbs, utils.CreateAccount, body, conf.ADMIN_EMAIL)
	})

	t.Run("valid_body_irrelevant_data_201_created", func(t *testing.T) {
		validBodyIrrelevantData := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":"%s","password":"%s","password_2":"%s","abc":123}`,
			conf.ADMIN_NAME, conf.ADMIN_EMAIL, conf.ADMIN_PHONE, helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountSuccess(
			t, app, dbs, utils.CreateAccount, validBodyIrrelevantData, conf.ADMIN_EMAIL,
		)
	})
}

func testCreateAccountClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, clientOperation, body string,
	expectedStatus int, expectedDetail string, expectedFieldErrors map[string][]string,
	expectedNonFieldErrors []string, expectedLog *models.Log,
) {
	resp := newRequestCreateAccount(t, app, clientOperation, body)
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

func testCreateAccountSuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, clientOperation, body, email string,
) {
	setup.SetUpApiGateway(t, dbs)

	var userCount int64
	helpers.CountUsers(t, dbs.ApiGateway, &userCount)
	require.EqualValues(t, 0, userCount)

	var sessionCount int64
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 0, sessionCount)

	resp := newRequestCreateAccount(t, app, clientOperation, body)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	helpers.CountUsers(t, dbs.ApiGateway, &userCount)
	require.EqualValues(t, 1, userCount)
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 1, sessionCount)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var createAcctRespBody controllers.CreateAccountResponseBody

		if err := json.Unmarshal(respBody, &createAcctRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Regexp(t, utils.TokenRegexp, createAcctRespBody.Token)

		var user models.User
		helpers.QueryTestUserByEmail(t, dbs.ApiGateway, &user, email)
		require.Equal(t, user.Slug, createAcctRespBody.User.Slug)
		require.Equal(t, user.Name, createAcctRespBody.User.Name)

		var session models.ClientSession
		helpers.QueryTestClientSessionLatest(t, dbs.ApiGateway, &session)
		require.Equal(t, session.TokenKey, createAcctRespBody.Token[:16])
		require.Equal(t, session.Digest, utils.HashToken(createAcctRespBody.Token))
	}
}

func newRequestCreateAccount(
	t *testing.T, app *fiber.App, clientOperation, body string,
) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/create_account", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", clientOperation)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
