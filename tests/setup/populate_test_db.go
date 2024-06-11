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
	validSessionTokens []string,
	expiredSessionTokens []string,
	validEmailTokens []string,
	expiredEmailTokens []string,
) {
	users = createTestUsers(t, dbs)
	validSessionTokens = createValidTestClientSessions(&users, t, dbs, conf)
	expiredSessionTokens = createExpiredTestClientSessions(&users, t, dbs, conf)
	validEmailTokens = createValidTestEmailTokens(&users, t, dbs)
	expiredEmailTokens = createExpiredTestEmailTokens(&users, t, dbs)
	return
}
