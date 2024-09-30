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

func CreateValidTestMFATokens(
	user *models.User, t *testing.T, dbs *databases.Databases,
) (tokens []controllers.AuthSecondFactorRequestBody) {

	var validMFATokens []models.MFAToken

	for i := 0; i < 2; i++ {
		if mfaTokenString, err := utils.GenerateSlug(80); err != nil {
			t.Fatalf("Generate test mfa token failed: %s", err.Error())
			panic(err)
		} else if oneTimePasscode, err := utils.GenerateOTP(); err != nil {
			t.Fatalf("Generate test otp failed: %s", err.Error())
			panic(err)
		} else {
			now := time.Now().UTC()

			validMFATokens = append(validMFATokens, models.MFAToken{
				UserSlug:  user.Slug,
				KeyDigest: utils.HashToken(mfaTokenString),
				OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
				TokenKey:  mfaTokenString[:16],
				CreatedAt: now.Add(time.Duration(1) * -time.Minute),
				ExpiresAt: now.Add(time.Duration(4) * time.Minute),
			})

			tokens = append(tokens, controllers.AuthSecondFactorRequestBody{
				MFAToken: mfaTokenString,
				PhoneOTP: strings.Join(oneTimePasscode, "")},
			)
		}
	}

	if result := dbs.ApiGateway.Create(&validMFATokens); result.Error != nil {
		t.Fatalf("Create test mfa tokens failed: %s", result.Error.Error())
		panic(result.Error)
	}

	return
}
