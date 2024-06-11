package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VerifyEmailTry(c *fiber.Ctx) error {
	var user models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(models.User); !ok {
		H.logger(c, utils.VerifyEmailTry, "UserContext not ok", "", "error", utils.ErrorUserContext)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if user.EmailIsVerified {
		H.logger(
			c, utils.VerifyEmailTry, "user.EmailIsVerified == true", user.EmailAddress, "warn",
			utils.ErrorBadClient,
		)

		return c.SendStatus(fiber.StatusNoContent)
	}

	var emailTokens []models.EmailVerificationToken

	if result := H.DBs.ApiGateway.Where("user_slug = ?", user.Slug).Order("expires_at DESC").
	Find(&emailTokens); result.Error != nil {
		H.logger(c, utils.VerifyEmailTry, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	var foundValidToken bool
	var invalidTokenKeys []string

	now := time.Now().UTC()
	t := now.Add(time.Duration(24) * time.Hour)
	
	// Clean up old tokens
	for _, token := range emailTokens {
		if !foundValidToken && token.ExpiresAt.After(now) && t.Sub(token.ExpiresAt).Seconds() <= 60 {
			foundValidToken = true
		} else {
			invalidTokenKeys = append(invalidTokenKeys, token.TokenKey)
		}
	}

	if result := H.DBs.ApiGateway.Where("token_key IN ?", invalidTokenKeys).
	Delete(&models.EmailVerificationToken{}); result.Error != nil {
		H.logger(c, utils.VerifyEmailTry, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if !foundValidToken {
		if tokenString, err := utils.GenerateSlug(80); err != nil {
			H.logger(c, utils.VerifyEmailTry, err.Error(), "", "error", "Failed generate email token")

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		} else if result := H.DBs.ApiGateway.Create(&models.EmailVerificationToken{
			UserSlug:  user.Slug,
			Digest:    utils.HashToken(tokenString),
			TokenKey:  tokenString[:16],
			CreatedAt: now,
			ExpiresAt: t,
		}); result.Error != nil {
			H.logger(
				c, utils.VerifyEmailTry, result.Error.Error(), "", "error", "Failed create email token",
			)

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		} else {
			H.sendEmail(user.EmailAddress, user.Name, tokenString)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
