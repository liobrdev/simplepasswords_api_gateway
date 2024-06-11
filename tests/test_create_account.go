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

	if conf.GO_FIBER_BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"name":"%s","email":"%s","password":"%s","password_2":"%s"}`

	setup.SetUpLogger(t, dbs)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testCreateAccountClientError(
			t, app, dbs, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			t, app, dbs, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testCreateAccountClientError(
			t, app, dbs, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(
			bodyFmt, "JDoe", "jdoe@email.co", helpers.VALID_PW, helpers.VALID_PW,
		) + `]`

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			t, app, dbs, "null", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			t, app, dbs, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testCreateAccountClientError(
			t, app, dbs, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			t, app, dbs, `"Valid JSON, but not an object."`, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
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
			t, app, dbs, "{}", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("missing_name_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"nmae":"%s","email":"%s","password":"%s","password_2":"%s"}`,
			"Spelled wrong!", "jdoe@email.co", helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			`{"name":null,"email":"%s","password":"%s","password_2":"%s"}`,
			"jdoe@email.co", helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		body := fmt.Sprintf(bodyFmt, "", "jdoe@email.co", helpers.VALID_PW, helpers.VALID_PW)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorAcctName,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_long_name_400_bad_request", func(t *testing.T) {
		if name, err := utils.GenerateSlug(65); err != nil {
			t.Fatalf("Generate long string failed: %s", err.Error())
		} else {
			body := fmt.Sprintf(
				bodyFmt, name, "jdoe@email.co", helpers.VALID_PW, helpers.VALID_PW,
			)

			testCreateAccountClientError(
				t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.CreateAccount,
					Detail:          name,
					Level:           "warn",
					Message:         utils.ErrorAcctName,
					RequestBody:     body,
				},
			)
		}
	})

	t.Run("missing_email_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","emial":"%s","password":"%s","password_2":"%s"}`,
			"JDoe", "spelled@wrong.co", helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			`{"name":"%s","email":null,"password":"%s","password_2":"%s"}`,
			"JDoe", helpers.VALID_PW, helpers.VALID_PW,
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		body := fmt.Sprintf(bodyFmt, "JDoe", "", helpers.VALID_PW, helpers.VALID_PW)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
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
			body := fmt.Sprintf(bodyFmt, "JDoe", email, helpers.VALID_PW, helpers.VALID_PW)

			testCreateAccountClientError(
				t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.CreateAccount,
					Detail:          email,
					Level:           "warn",
					Message:         utils.ErrorAcctEmail,
					RequestBody:     body,
				},
			)
		}
	})

	t.Run("missing_password_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(
			`{"name":"%s","email":"%s","passwrod":"%s","password_2":"%s"}`,
			"JDoe", "jdoe@email.co", "$pelledWr0ng1234", "$pelledWr0ng1234",
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
			bodyFmt, "JDoe", "jdoe@email.co", helpers.VALID_PW, "N0nmatchpa$sw0rd",
		)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", "jdoe@email.co", "jdoe@email.co")

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		if slug, err := utils.GenerateSlug(253); err != nil {
			t.Fatalf("Generate long password failed: %s", err.Error())
		} else {
			password := slug + "qQ1!"
			body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

			testCreateAccountClientError(
				t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.CreateAccount,
					Detail:          "Too long: 257 > 256",
					Level:           "warn",
					Message:         utils.ErrorAcctPW,
					RequestBody:     body,
				},
			)
		}
	})

	t.Run("too_short_password_400_bad_request", func(t *testing.T) {
		if slug, err := utils.GenerateSlug(11); err != nil {
			t.Fatalf("Generate short password failed: %s", err.Error())
		} else {
			password := slug + "qQ1!"
			body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

			testCreateAccountClientError(
				t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
					ClientIP:        clientIP,
					ClientOperation: utils.CreateAccount,
					Detail:          "Too short: 15 < 16",
					Level:           "warn",
					Message:         utils.ErrorAcctPW,
					RequestBody:     body,
				},
			)
		}
	})

	t.Run("password_missing_uppercase_400_bad_request", func(t *testing.T) {
		password := "n0uppercasepa$sw0rd"
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		password := "N0LOWERCASEPA$SW0RD"
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		password := "nONumberPA$SWoRD"
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		password := "nOSp3cialCharPA5SWoRD"
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
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
		password := "h@s Spa(3	PA5SWoRD"
		body := fmt.Sprintf(bodyFmt, "JDoe", "jdoe@email.co", password, password)

		testCreateAccountClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.CreateAccount,
				Detail:          "Has whitespace",
				Level:           "warn",
				Message:         utils.ErrorAcctPW,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_email_already_exists_409_conflict", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		users, _, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		email := users[0].EmailAddress
		body := fmt.Sprintf(bodyFmt, "JDoe", email, helpers.VALID_PW, helpers.VALID_PW)

		testCreateAccountClientError(
			t, app, dbs, body, http.StatusConflict, utils.ErrorDiffEmail, nil, nil, nil,
		)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_201_created", func(t *testing.T) {
		email := "jake.doe@email.co"
		body := fmt.Sprintf(bodyFmt, "Jake Doe", email, helpers.VALID_PW, helpers.VALID_PW)
		testCreateAccountSuccess(t, app, dbs, body, email)
	})

	t.Run("valid_body_irrelevant_data_201_created", func(t *testing.T) {
		email := "jake.doe@email.co"
		validBodyIrrelevantData := fmt.Sprintf(
			`{"name":"%s","email":"%s","password":"%s","password_2":"%s","abc":123}`,
			"Jake Doe", email, helpers.VALID_PW, helpers.VALID_PW,
		)
		testCreateAccountSuccess(t, app, dbs, validBodyIrrelevantData, email)
	})
}

func testCreateAccountClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestCreateAccount(t, app, body)
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
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, email string,
) {
	setup.SetUpApiGateway(t, dbs)

	var userCount int64
	helpers.CountUsers(t, dbs.ApiGateway, &userCount)
	require.EqualValues(t, 0, userCount)

	var sessionCount int64
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 0, sessionCount)

	resp := newRequestCreateAccount(t, app, body)
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
		require.Equal(t, user.EmailAddress, createAcctRespBody.User.EmailAddress)
		require.Equal(t, user.EmailIsVerified, createAcctRespBody.User.EmailIsVerified)
		require.Equal(t, false, createAcctRespBody.User.EmailIsVerified)
		require.Equal(t, user.PhoneNumber, createAcctRespBody.User.PhoneNumber)
		require.Equal(t, "", createAcctRespBody.User.PhoneNumber)
		require.Equal(t, user.PhoneIsVerified, createAcctRespBody.User.PhoneIsVerified)
		require.Equal(t, false, createAcctRespBody.User.PhoneIsVerified)
		require.Equal(t, user.MfaIsEnabled, createAcctRespBody.User.MfaIsEnabled)
		require.Equal(t, false, createAcctRespBody.User.MfaIsEnabled)
		require.Equal(t, user.IsActive, createAcctRespBody.User.IsActive)
		require.Equal(t, true, createAcctRespBody.User.IsActive)
		require.Empty(t, createAcctRespBody.User.CreatedAt)
		require.Empty(t, createAcctRespBody.User.UpdatedAt)
		require.Empty(t, createAcctRespBody.User.PasswordSalt)
		require.Empty(t, createAcctRespBody.User.PasswordHash)

		var session models.ClientSession
		helpers.QueryTestClientSessionLatest(t, dbs.ApiGateway, &session)
		require.Equal(t, session.UserSlug, createAcctRespBody.User.Slug)
		require.Equal(t, session.TokenKey, createAcctRespBody.Token[:16])
		require.Equal(t, session.Digest, utils.HashToken(createAcctRespBody.Token))
	}
}

func newRequestCreateAccount(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest(http.MethodPost, "/api/auth/create_account", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.CreateAccount)

	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
