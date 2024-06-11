package setup

import (
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func createExpiredTestEmailTokens(
	users *[]models.User, t *testing.T, dbs *databases.Databases,
) (tokens []string) {

	var expiredTokens []models.EmailVerificationToken
	now := time.Now().UTC()

	for _, user := range *users {
		for j := 0; j < 2; j++ {
			if token, err := utils.GenerateSlug(80); err != nil {
				t.Fatalf("Generate test email token failed: %s", err.Error())
				panic(err)
			} else {
				expiredTokens = append(expiredTokens, models.EmailVerificationToken{
					UserSlug:  user.Slug,
					Digest:    utils.HashToken(token),
					TokenKey:  token[:16],
					CreatedAt: now.Add(time.Duration(-25) * time.Hour),
					ExpiresAt: now.Add(time.Duration(-1) * time.Hour),
				})

				tokens = append(tokens, token)
			}
		}
	}

	if result := dbs.ApiGateway.Create(&expiredTokens); result.Error != nil {
		t.Fatalf("Create test email tokens failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
