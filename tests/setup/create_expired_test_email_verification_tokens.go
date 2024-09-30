package setup

import (
	"strings"
	"testing"
	"time"

	"github.com/liobrdev/simplepasswords_api_gateway/controllers"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func CreateExpiredTestEmailVerificationTokens(
	user *models.User, t *testing.T, dbs *databases.Databases,
) (tokens []controllers.VerifyEmailTryResponseBody) {

	var expiredEmailTokens []models.EmailVerificationToken

	for i := 0; i < 2; i++ {
		if tokenString, err := utils.GenerateSlug(80); err != nil {
			t.Fatalf("Generate test email token failed: %s", err.Error())
			panic(err)
		} else if oneTimePasscode, err := utils.GenerateOTP(); err != nil {
			t.Fatalf("Generate test otp failed: %s", err.Error())
			panic(err)
		} else {
			now := time.Now().UTC()

			expiredEmailTokens = append(expiredEmailTokens, models.EmailVerificationToken{
				UserSlug:  user.Slug,
				KeyDigest: utils.HashToken(tokenString),
				OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
				TokenKey:  tokenString[:16],
				CreatedAt: now.Add(time.Duration(11) * -time.Minute),
				ExpiresAt: now.Add(time.Duration(1) * -time.Minute),
			})

			tokens = append(tokens, controllers.VerifyEmailTryResponseBody{
				TestEmailToken:	tokenString,
				TestOTP: 				strings.Join(oneTimePasscode, "")},
			)
		}
	}

	if result := dbs.ApiGateway.Create(&expiredEmailTokens); result.Error != nil {
		t.Fatalf("Create test email tokens failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
