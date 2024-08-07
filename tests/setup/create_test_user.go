package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createTestUser(t *testing.T, dbs *databases.Databases, conf *config.AppConfig) models.User {
	if hash, err := utils.GenerateUserCredentials(helpers.VALID_PW, conf); err != nil {
		t.Fatalf("Generate test user credentials failed: %s", err.Error())
		panic(err)
	} else {
		user := models.User{
			Slug:         helpers.NewSlug(t),
			Name:         conf.ADMIN_NAME,
			EmailAddress: conf.ADMIN_EMAIL,
			PhoneNumber: 	conf.ADMIN_PHONE,
			PasswordHash: hash,
		}

		if result := dbs.ApiGateway.Create(&user); result.Error != nil {
			t.Fatalf("Create test user failed: %s", result.Error.Error())
			panic(result.Error)
		}

		return user
	}
}
