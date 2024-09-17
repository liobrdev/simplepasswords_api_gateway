package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/controllers"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func populateTestDBApiGateway(t *testing.T, dbs *databases.Databases, conf *config.AppConfig) (
	user models.User,
	validSessionTokens, expiredSessionTokens []string,
	validMFATokens, expiredMFATokens []controllers.AuthSecondFactorRequestBody,
) {
	user = createTestUser(t, dbs)
	validSessionTokens = createValidTestClientSessions(&user, t, dbs, conf)
	expiredSessionTokens = createExpiredTestClientSessions(&user, t, dbs, conf)
	validMFATokens = createValidTestMFATokens(&user, t, dbs)
	expiredMFATokens = createExpiredTestMFATokens(&user, t, dbs)
	return
}
