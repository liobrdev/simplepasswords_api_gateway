package helpers

import (
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func LogInTestUsers(
	users *[]models.User,
	t *testing.T,
	dbs *databases.Databases,
) (tokens []string) {
	var sessions []models.ClientSession
	createdAt := time.Now().UTC()

	for i := 0; i < len(*users); i++ {
		if token, err := utils.GenerateSlug(80); err != nil {
			t.Fatalf("Generate test client session token failed: %s", err.Error())
			panic(err)
		} else {
			sessions = append(sessions, models.ClientSession{
				UserSlug:  (*users)[i].Slug,
				Digest:    utils.HashToken(token),
				TokenKey:  token[:16],
				CreatedAt: createdAt,
				ExpiresAt: createdAt.Add(time.Duration(15) * time.Minute),
			})

			tokens = append(tokens, token)
		}
	}

	if result := dbs.ApiGateway.Create(&sessions); result.Error != nil {
		t.Fatalf("Create test client sessions failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
