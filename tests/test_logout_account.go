package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/setup"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func testLogoutAccount(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	t.Run("valid_logout_204_no_content", func(t *testing.T) {
		user := setup.SetUpApiGatewayWithData(t, dbs)
		validTokens := setup.CreateValidTestClientSessions(&user, t, dbs, conf)

		var sessionCount int64
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 2, sessionCount)

		resp := newRequestLogoutAccount(t, app, "Token " + validTokens[0])
		require.Equal(t, http.StatusNoContent, resp.StatusCode)

		if respBody, err := io.ReadAll(resp.Body); err != nil {
			t.Fatalf("Read response body failed: %s", err.Error())
		} else {
			require.Empty(t, respBody)
		}

		var session models.ClientSession

		if result := dbs.ApiGateway.First(&session, "token_key = ?", validTokens[0][:16]);
		result.Error != nil {
			require.ErrorIs(t, result.Error, gorm.ErrRecordNotFound)
		} else {
			t.Fatalf("Deleted client session query failed: %s", result.Error.Error())
		}

		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 1, sessionCount)
	})
}

func newRequestLogoutAccount(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout_account", nil)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Client-Operation", utils.LogoutAccount)
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
