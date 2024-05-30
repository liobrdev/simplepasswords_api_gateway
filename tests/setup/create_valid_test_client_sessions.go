package setup

import (
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createValidTestClientSessions(
	users *[]models.User, t *testing.T, dbs *databases.Databases, conf *config.AppConfig,
) (tokens []string) {

	var validSessions []models.ClientSession
	now := time.Now().UTC()

	for _, user := range *users {
		for j, clientIP := 0, "0.0.0.0"; j < 2; j++ {
			if j == 1 {
				clientIP = helpers.OLD_IP
			} else if conf.GO_FIBER_BEHIND_PROXY {
				clientIP = helpers.CLIENT_IP
			}

			if token, err := utils.GenerateSlug(80); err != nil {
				t.Fatalf("Generate test client session token failed: %s", err.Error())
				panic(err)
			} else {
				validSessions = append(validSessions, models.ClientSession{
					UserSlug:  user.Slug,
					ClientIP:  clientIP,
					Digest:    utils.HashToken(token),
					TokenKey:  token[:16],
					CreatedAt: now,
					ExpiresAt: now.Add(time.Duration(15) * time.Minute),
				})

				tokens = append(tokens, token)
			}
		}
	}

	if result := dbs.ApiGateway.Create(&validSessions); result.Error != nil {
		t.Fatalf("Create test client sessions failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
