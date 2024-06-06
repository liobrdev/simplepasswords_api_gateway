package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func populateTestDBApiGateway(
	t *testing.T, dbs *databases.Databases, conf *config.AppConfig,
) (
	users []models.User,
	validTokens []string,
	expiredTokens []string,
) {
	users = createTestUsers(t, dbs)
	validTokens = createValidTestClientSessions(&users, t, dbs, conf)
	expiredTokens = createExpiredTestClientSessions(&users, t, dbs, conf)
	return
}
