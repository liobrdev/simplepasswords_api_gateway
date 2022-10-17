package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func populateTestDBApiGateway(t *testing.T, dbs *databases.Databases) (
	*[]models.User,
	*[]models.DeactivatedUser,
) {
	users := createTestUsers(t, dbs)
	deactivatedUsers := createTestDeactivatedUsers(t, dbs)
	return &users, &deactivatedUsers
}
