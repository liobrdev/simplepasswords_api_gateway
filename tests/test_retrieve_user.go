package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/setup"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func testRetrieveUser(
	t *testing.T, app *fiber.App, dbs *databases.Databases, conf *config.AppConfig,
) {
	t.Run("valid_token_200_ok", func(t *testing.T) {
		user, validSessionTokens, _, _, _ := setup.SetUpApiGatewayWithData(t, dbs, conf)

		var sessionCount int64
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 4, sessionCount)

		testRetrieveUserSuccess(t, app, user.Slug, user.Name, "Token " + validSessionTokens[0])
		helpers.CountClientSessions(t, dbs.ApiGateway, &sessionCount)
		require.EqualValues(t, 2, sessionCount)
	})
}

func testRetrieveUserSuccess(t *testing.T, app *fiber.App, slug, name, authHeader string) {
	resp := newRequestRetrieveUser(t, app, authHeader)
	require.Equal(t, 200, resp.StatusCode)

	if respBody, err := io.ReadAll(resp.Body); err != nil {
		t.Fatalf("Read response body failed: %s", err.Error())
	} else {
		var user models.User

		if err := json.Unmarshal(respBody, &user); err != nil {
			t.Fatalf("JSON unmarshal failed: %s", err.Error())
		}

		require.Equal(t, slug, user.Slug)
		require.Equal(t, name, user.Name)
	}
}

func newRequestRetrieveUser(t *testing.T, app *fiber.App, authHeader string) *http.Response {
	req := httptest.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", helpers.CLIENT_IP)
	req.Header.Set("Client-Operation", utils.RetrieveUser)
	req.Header.Set("Authorization", authHeader)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Send test request failed: %s", err.Error())
	}

	return resp
}
