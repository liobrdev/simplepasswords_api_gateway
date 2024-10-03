package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type VerifyEmailTryResponseBody struct {
	TestEmailToken string `json:"test_email_token,omitempty"`
	TestOTP	 			 string `json:"test_otp,omitempty"`
}

func (H Handler) VerifyEmailTry(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.VerifyEmailTry {
		H.logger(c, utils.VerifyEmailTry, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var session *models.ClientSession
	var ok bool

	if session, ok = c.UserContext().Value(sessionContextKey{}).(*models.ClientSession); !ok {
		H.logger(c, utils.VerifyEmailTry, "", "", "error", "Failed session.User context", "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	user := &session.User

	if user.EmailIsVerified {
		H.logger(c, utils.VerifyEmailTry, "", "", "warn", utils.ErrorAlreadyVerified, user.Slug)

		return c.Status(200).JSON(&VerifyEmailTryResponseBody{})
	}

	var emailTokens []models.EmailVerificationToken

	// Get all user's email verification tokens for cleanup
	if result := H.DBs.ApiGateway.Where("user_slug = ?", user.Slug).Find(&emailTokens);
	result.Error != nil {
		H.logger(
			c, utils.VerifyEmailTry, result.Error.Error(), "", "error", utils.ErrorFailedDB, user.Slug,
		)
	}

	var currentToken *models.EmailVerificationToken
	now := time.Now().UTC()

	// Clean up expired email tokens
	for _, token := range emailTokens {
		if !token.ExpiresAt.After(now) {
			H.deleteEmailToken(c, &token)
		} else if currentToken == nil {
			currentToken = &token
		} else if !token.ExpiresAt.Before(currentToken.ExpiresAt) {
			H.deleteEmailToken(c, currentToken)
			currentToken = &token
		} else {
			H.deleteEmailToken(c, &token)
		}
	}

	if currentToken != nil {
		later := now.Add(time.Duration(10) * time.Minute)

		if later.Sub(currentToken.ExpiresAt).Seconds() < 30 {
			H.logger(c, utils.VerifyEmailTry, "", "", "warn", "Too soon retry", user.Slug)

			return c.Status(200).JSON(&VerifyEmailTryResponseBody{})
		} else {
			H.deleteEmailToken(c, currentToken)
		}
	}

	var tokenString string
	var oneTimePasscode []string
	var err error

	if tokenString, err = utils.GenerateSlug(80); err != nil {
		H.logger(c, utils.VerifyEmailTry, err.Error(), "", "error", "Failed generate string", user.Slug)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if oneTimePasscode, err = utils.GenerateOTP(); err != nil {
		H.logger(c, utils.VerifyEmailTry, err.Error(), "", "error", "Failed generate otp", user.Slug)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if result := H.DBs.ApiGateway.Create(&models.EmailVerificationToken{
		UserSlug:  user.Slug,
		KeyDigest: utils.HashToken(tokenString),
		OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
		TokenKey:	 tokenString[:16],
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(10) * time.Minute),
	}); result.Error != nil {
		H.logger(
			c, utils.VerifyEmailTry, result.Error.Error(), "", "error", "Failed create token", user.Slug,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	var testEmailToken string
	var testOTP string

	if H.Conf.ENVIRONMENT == "testing" {
		testEmailToken = tokenString
		testOTP = strings.Join(oneTimePasscode, "")
	} else if err = H.sendVerificationEmail(user, tokenString, oneTimePasscode); err != nil {
		H.logger(
			c, utils.VerifyEmailTry, err.Error(), "", "error", "Failed send verification email",
			user.Slug,
		)

		if result := H.DBs.ApiGateway.Exec(
			"DELETE FROM email_verification_tokens WHERE token_key = ?", tokenString[:16],
		); result.Error != nil {
			H.logger(
				c, utils.VerifyEmailTry, result.Error.Error(), "", "error", "Failed delete token",
				user.Slug,
			)
		}

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).JSON(&VerifyEmailTryResponseBody{
		TestEmailToken: testEmailToken,
		TestOTP: 	testOTP,
	})
}

func (H Handler) deleteEmailToken(c *fiber.Ctx, token *models.EmailVerificationToken) {
	if result := H.DBs.ApiGateway.Delete(&token); result.Error != nil {
		H.logger(
			c, utils.VerifyEmailTry, result.Error.Error(), "", "error", utils.ErrorFailedDB,
			token.UserSlug,
		)
	} else if n := result.RowsAffected; n != 1 {
		H.logger(
			c, utils.VerifyEmailTry, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB, token.UserSlug,
		)
	}
}

func (H Handler) sendVerificationEmail(user *models.User, token string, otp []string) error {
	return H.sendEmail(
		"Verify your email address", H.Conf.SUPPORT_EMAIL, []string{user.EmailAddress},
		"email_verify_email.html", map[string]string{
			"Name": user.Name,
			"Otp": strings.Join(otp, " "),
			"Link": H.Conf.APP_SCHEME + "://" + H.Conf.APP_DOMAIN + "/verify_email?token=" + token,
		},
	)
}
