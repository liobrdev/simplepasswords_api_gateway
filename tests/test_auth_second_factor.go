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

func testAuthSecondFactor(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	var clientIP string

	if conf.BEHIND_PROXY {
		clientIP = helpers.CLIENT_IP
	} else {
		clientIP = "0.0.0.0"
	}

	bodyFmt := `{"mfa_token":"%s","phone_otp":"%s"}`

	var dummyMFAToken string
	var dummyPhoneOTP string
	var err error

	if dummyMFAToken, err = utils.GenerateSlug(80); err != nil {
		t.Fatalf("Generate dummy MFA token failed: %s", err.Error())
	}

	if blocks, err := utils.GenerateOTP(); err != nil {
		t.Fatalf("Generate dummy phone OTP failed: %s", err.Error())
	} else {
		dummyPhoneOTP = strings.Join(blocks, "")
	}

	setup.SetUpLogger(t, dbs)

	t.Run("empty_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, "", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character '\x00' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "",
			},
		)
	})

	t.Run("array_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, "[]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[]",
			},
		)

		testAuthSecondFactorClientError(
			t, app, dbs, "[{}]", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "[{}]",
			},
		)

		body := `[` + fmt.Sprintf(bodyFmt, dummyMFAToken, dummyPhoneOTP) + `]`

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character '[' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, "null", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     "null",
			},
		)
	})

	t.Run("boolean_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, "true", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character 't' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "true",
			},
		)

		testAuthSecondFactorClientError(
			t, app, dbs, "false", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "invalid character 'f' looking for beginning of value",
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     "false",
			},
		)
	})

	t.Run("string_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, `"Valid JSON, but not an object."`, 400, utils.ErrorBadRequest, nil, nil,
			&models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          `invalid character '"' looking for beginning of value`,
				Level:           "warn",
				Message:         utils.ErrorParse,
				RequestBody:     `"Valid JSON, but not an object."`,
			},
		)
	})

	t.Run("empty_object_body_400_bad_request", func(t *testing.T) {
		testAuthSecondFactorClientError(
			t, app, dbs, "{}", 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     "{}",
			},
		)
	})

	t.Run("missing_mfa_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(`{"maf_token":"%s","phone_otp":"%s"}`, dummyMFAToken, dummyPhoneOTP)

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_mfa_token_400_bad_request", func(t *testing.T) {
		body := `{"mfa_token":null,"phone_otp":"` + dummyPhoneOTP + `"}`

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_mfa_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, "", dummyPhoneOTP)

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_long_mfa_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, dummyMFAToken + "1", dummyPhoneOTP)

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          dummyMFAToken + "1",
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_short_mfa_token_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, dummyMFAToken[:79], dummyPhoneOTP)

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          dummyMFAToken[:79],
				Level:           "warn",
				Message:         utils.ErrorMFAToken,
				RequestBody:     body,
			},
		)
	})

	t.Run("missing_phone_otp_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(`{"mfa_token":"%s","phoen_otp":"%s"}`, dummyMFAToken, dummyPhoneOTP)

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorPhoneOTP,
				RequestBody:     body,
			},
		)
	})

	t.Run("null_phone_otp_400_bad_request", func(t *testing.T) {
		body := `{"mfa_token":"` + dummyMFAToken + `","phone_otp":null}`

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorPhoneOTP,
				RequestBody:     body,
			},
		)
	})

	t.Run("empty_phone_otp_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, dummyMFAToken, "")

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          "",
				Level:           "warn",
				Message:         utils.ErrorPhoneOTP,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_long_phone_otp_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, dummyMFAToken, dummyPhoneOTP + "1")

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          dummyPhoneOTP + "1",
				Level:           "warn",
				Message:         utils.ErrorPhoneOTP,
				RequestBody:     body,
			},
		)
	})

	t.Run("too_short_phone_otp_400_bad_request", func(t *testing.T) {
		body := fmt.Sprintf(bodyFmt, dummyMFAToken, dummyPhoneOTP[:19])

		testAuthSecondFactorClientError(
			t, app, dbs, body, 400, utils.ErrorBadRequest, nil, nil, &models.Log{
				ClientIP:        clientIP,
				ClientOperation: utils.AuthSecondFactor,
				Detail:          dummyPhoneOTP[:19],
				Level:           "warn",
				Message:         utils.ErrorPhoneOTP,
				RequestBody:     body,
			},
		)
	})

	t.Run("valid_body_failed_authenticate_400_bad_request", func(t *testing.T) {
		setup.SetUpLogger(t, dbs)
		_, _, _, validMFATokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		body := fmt.Sprintf(bodyFmt, dummyMFAToken, validMFATokens[0].PhoneOTP)
		testAuthSecondFactorClientError(t, app, dbs, body, 400, utils.ErrorAuthenticate, nil, nil, nil)

		body = fmt.Sprintf(bodyFmt, validMFATokens[0].MFAToken, dummyPhoneOTP)
		testAuthSecondFactorClientError(t, app, dbs, body, 400, utils.ErrorAuthenticate, nil, nil, nil)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_token_expired_400_bad request", func(t *testing.T) {
		_, _, _, _, expiredMFATokens := setup.SetUpApiGatewayWithData(t, dbs, conf)

		body := fmt.Sprintf(bodyFmt, expiredMFATokens[1].MFAToken, expiredMFATokens[1].PhoneOTP)
		testAuthSecondFactorClientError(t, app, dbs, body, 400, utils.ErrorAuthenticate, nil, nil, nil)

		var logCount int64
		helpers.CountLogs(t, dbs.Logger, &logCount)
		require.EqualValues(t, 0, logCount)
	})

	t.Run("valid_body_200_ok", func(t *testing.T) {
		_, _, _, validMFATokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		body := fmt.Sprintf(bodyFmt, validMFATokens[0].MFAToken, validMFATokens[0].PhoneOTP)
		testAuthSecondFactorSuccess(t, app, dbs, body, conf.ADMIN_EMAIL)
	})

	t.Run("valid_body_irrelevant_data_200_ok", func(t *testing.T) {
		_, _, _, validMFATokens, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)
		validBodyIrrelevantData := fmt.Sprintf(
			`{"mfa_token":"%s","phone_otp":"%s","abc":123}`,
			validMFATokens[0].MFAToken, validMFATokens[0].PhoneOTP,
		)
		testAuthSecondFactorSuccess(t, app, dbs, validBodyIrrelevantData, conf.ADMIN_EMAIL)
	})
}

func testAuthSecondFactorClientError(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body string, expectedStatus int,
	expectedDetail string, expectedFieldErrors map[string][]string, expectedNonFieldErrors []string,
	expectedLog *models.Log,
) {
	resp := newRequestAuthSecondFactor(t, app, body)
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

func testAuthSecondFactorSuccess(
	t *testing.T, app *fiber.App, dbs *databases.Databases, body, email string,
) {
	var sessionCount int64
	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 4, sessionCount)

	resp := newRequestAuthSecondFactor(t, app, body)
	require.Equal(t, 200, resp.StatusCode)

	helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
	require.EqualValues(t, 5, sessionCount)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var authSecondFactorRespBody controllers.AuthSecondFactorResponseBody

		if err := json.Unmarshal(respBody, &authSecondFactorRespBody); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Regexp(t, utils.TokenRegexp, authSecondFactorRespBody.Token)

		var user models.User
		helpers.QueryTestUserByEmail(t, dbs.ApiGateway, &user, email)
		require.Equal(t, user.Name, authSecondFactorRespBody.UserName)		

		var session models.ClientSession
		helpers.QueryTestClientSessionLatest(t, dbs.ApiGateway, &session)
		require.Equal(t, session.TokenKey, authSecondFactorRespBody.Token[:16])
		require.Equal(t, session.Digest, utils.HashToken(authSecondFactorRespBody.Token))
	}
}

func newRequestAuthSecondFactor(t *testing.T, app *fiber.App, body string) *http.Response {
	reqBody := strings.NewReader(body)
	req := httptest.NewRequest("POST", "/api/auth/second_factor", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.AuthSecondFactor)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
