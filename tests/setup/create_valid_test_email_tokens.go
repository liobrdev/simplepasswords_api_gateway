package setup

import (
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createValidTestEmailTokens(
	users *[]models.User, t *testing.T, dbs *databases.Databases,
) (tokens []string) {
	var validTokens []models.EmailVerificationToken

	for _, user := range *users {
		for j := 0; j < 2; j++ {
			if token, err := utils.GenerateSlug(80); err != nil {
				t.Fatalf("Generate test email token failed: %s", err.Error())
				panic(err)
			} else {
				now := time.Now().UTC()

				validTokens = append(validTokens, models.EmailVerificationToken{
					UserSlug:  user.Slug,
					Digest:    utils.HashToken(token),
					TokenKey:  token[:16],
					CreatedAt: now,
					ExpiresAt: now.Add(time.Duration(24) * time.Hour),
				})

				tokens = append(tokens, token)
			}
		}
	}

	if result := dbs.ApiGateway.Create(&validTokens); result.Error != nil {
		t.Fatalf("Create test email tokens failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
