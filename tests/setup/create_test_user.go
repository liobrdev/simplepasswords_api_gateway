package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createTestUser(t *testing.T, dbs *databases.Databases) models.User {
	password := utils.HashToken(helpers.VALID_EMAIL_1 + helpers.VALID_PW_1)

	if hash, salt, err := utils.GenerateUserCredentials(password); err != nil {
		t.Fatalf("Generate test user credentials failed: %s", err.Error())
		panic(err)
	} else {
		user := models.User{
			Slug:         helpers.NewSlug(t),
			Name:         helpers.VALID_NAME_1,
			EmailAddress: helpers.VALID_EMAIL_1,
			PhoneNumber: 	helpers.VALID_PHONE_1,
			PasswordHash: hash,
			PasswordSalt: salt,
		}

		if result := dbs.ApiGateway.Create(&user); result.Error != nil {
			t.Fatalf("Create test user failed: %s", result.Error.Error())
			panic(result.Error)
		}

		return user
	}
}
