package setup

import (
	"testing"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/tests/helpers"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createTestUsers(t *testing.T, dbs *databases.Databases) []models.User {
	if salt1, hash1, err := utils.GenerateUserCredentials(helpers.VALID_PW); err != nil {
		t.Fatalf("Generate test user credentials failed: %s", err.Error())
		panic(err)
	} else if salt2, hash2, err := utils.GenerateUserCredentials(helpers.VALID_PW); err != nil {
		t.Fatalf("Generate test user credentials failed: %s", err.Error())
		panic(err)
	} else {
		users := []models.User{
			{
				Slug:         helpers.NewSlug(t),
				Name:         "Jane Doe",
				EmailAddress: "jane.doe@email.co",
				PasswordSalt: salt1,
				PasswordHash: hash1,
			},
			{
				Slug:         helpers.NewSlug(t),
				Name:         "John Doe",
				EmailAddress: "john.doe@email.co",
				PasswordSalt: salt2,
				PasswordHash: hash2,
			},
		}

		if result := dbs.ApiGateway.Create(&users); result.Error != nil {
			t.Fatalf("Create test users failed: %s", result.Error.Error())
			panic(result.Error)
		}

		return users
	}
}
