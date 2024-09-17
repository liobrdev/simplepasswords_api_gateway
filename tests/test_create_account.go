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

	bodyFmt := `{"name":"%s","email":"%s","phone":"%s","password":"%s"}`
	
	setup.SetUpLogger(t, dbs)

	t.Run("wrong_client_operation_header_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
			bodyFmt, helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
			bodyFmt, "@()^!~", helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "@()^!~",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_name_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"nmae":"%s","email":"%s","phone":"%s","password":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
			`{"name":null,"email":"%s","phone":"%s","password":"%s"}`,
			helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
		body := fmt.Sprintf(bodyFmt, "", helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2)

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

	t.Run("invalid_email_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		body := fmt.Sprintf(
			bodyFmt, helpers.VALID_NAME_2, "!$@&^&.o", helpers.VALID_PHONE_2, helpers.HexHash2,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "!$@&^&.o",
				Level:           "warn",
				Message:         utils.ErrorAcctEmail,
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","emial":"%s","phone":"%s","password":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
			`{"name":"%s","email":null,"phone":"%s","password":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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
		body := fmt.Sprintf(bodyFmt, helpers.VALID_NAME_2, "", helpers.VALID_PHONE_2, helpers.HexHash2)

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

	t.Run("invalid_phone_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		body := fmt.Sprintf(bodyFmt, helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, "123", helpers.HexHash2)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "123",
				Level:           "warn",
				Message:         utils.ErrorAcctPhone,
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_phone_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","phnoe":"%s","password":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPhone,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_phone_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":null,"password":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.HexHash2,
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPhone,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_phone_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, "", helpers.HexHash2)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctPhone,
				RequestBody:     body,
			},
		)
	})

	t.Run("invalid_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, "0123456789abcdeg",
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "encoding/hex: invalid byte: U+0067 'g'",
				Level:           "error",
				Message:         "Failed decode password",
				RequestBody:     body,
			},
		)

		body = fmt.Sprintf(
			bodyFmt,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, "0123456789abcde",
		)

		testCreateAccountClientError(
			t, app, dbs, utils.CreateAccount, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "encoding/hex: odd length hex string",
				Level:           "error",
				Message:         "Failed decode password",
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":"%s","passwrod":"%s"}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
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

	t.Run("null_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":"%s","password":null}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2,
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

	t.Run("empty_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			bodyFmt, helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, "",
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

	t.Run("valid_body_user_already_exists_500_error", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		setup.SetUpApiGatewayWithData(t, dbs, conf)

		body := fmt.Sprintf(
			bodyFmt, helpers.VALID_NAME_1, helpers.VALID_EMAIL_1, helpers.VALID_PHONE_1, helpers.HexHash1,
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
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
		)

		testCreateAccountSuccess(t, app, dbs, utils.CreateAccount, body, helpers.VALID_EMAIL_2)
	})

	t.Run("valid_body_irrelevant_data_201_created", func(t *testing.T) {
		validBodyIrrelevantData := fmt.Sprintf(
			`{"name":"%s","email":"%s","phone":"%s","password":"%s","abc":123}`,
			helpers.VALID_NAME_2, helpers.VALID_EMAIL_2, helpers.VALID_PHONE_2, helpers.HexHash2,
		)

		testCreateAccountSuccess(
			t, app, dbs, utils.CreateAccount, validBodyIrrelevantData, helpers.VALID_EMAIL_2,
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
