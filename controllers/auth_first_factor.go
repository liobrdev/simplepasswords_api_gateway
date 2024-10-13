package controllers

import (
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type AuthFirstFactorRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthFirstFactorResponseBody struct {
	MFAToken string `json:"mfa_token"`
	TestOTP	 string `json:"test_otp,omitempty"`
}

func (H Handler) AuthFirstFactor(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.AuthFirstFactor {
		H.logger(c, utils.AuthFirstFactor, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	body := AuthFirstFactorRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "warn", utils.ErrorParse, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.EmailRegexp.Match([]byte(body.Email)) {
		H.logger(c, utils.AuthFirstFactor, body.Email, "", "warn", utils.ErrorAcctEmail, "")

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.AuthFirstFactor, "", "", "warn", utils.ErrorAcctPW, "")

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	var user models.User

	if result := H.DBs.ApiGateway.Where("email_address = ?", body.Email).Limit(1).Find(&user);
	result.Error != nil {
		H.logger(c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB, "")

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	} else if n != 1 {
		H.logger(
			c, utils.AuthFirstFactor, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB, "",
		)

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if !user.IsActive {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if password, err := hex.DecodeString(body.Password); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", utils.ErrorAcctPW, user.Slug)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	} else if !utils.CompareHashAndPassword(user.PasswordHash, password, user.PasswordSalt) {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	var mfaTokens []models.MFAToken

	// Get all user's mfa tokens for cleanup
	if result := H.DBs.ApiGateway.Where("user_slug = ?", user.Slug).Find(&mfaTokens);
	result.Error != nil {
		H.logger(
			c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB, user.Slug,
		)
	}

	var currentToken *models.MFAToken
	now := time.Now().UTC()

	// Clean up expired mfa tokens
	for _, token := range mfaTokens {
		if !token.ExpiresAt.After(now) {
			H.deleteMfaToken(c, &token)
		} else if currentToken == nil {
			currentToken = &token
		} else if !token.ExpiresAt.Before(currentToken.ExpiresAt) {
			H.deleteMfaToken(c, currentToken)
			currentToken = &token
		} else {
			H.deleteMfaToken(c, &token)
		}
	}

	if currentToken != nil {
		later := now.Add(time.Duration(5) * time.Minute)

		if later.Sub(currentToken.ExpiresAt).Seconds() < 30 {
			H.logger(c, utils.AuthFirstFactor, "", "", "warn", "Too soon retry", user.Slug)
			time.Sleep(time.Second)

			return c.Status(200).JSON(&AuthFirstFactorRequestBody{})
		} else {
			H.deleteMfaToken(c, currentToken)
		}
	}

	var tokenString	string
	var oneTimePasscode []string
	var err error

	if tokenString, err = utils.GenerateSlug(80); err != nil {
		H.logger(
			c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed generate mfa string", user.Slug,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if oneTimePasscode, err = utils.GenerateOTP(); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed generate otp", user.Slug)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if result := H.DBs.ApiGateway.Create(&models.MFAToken{
		UserSlug:  user.Slug,
		KeyDigest: utils.HashToken(tokenString),
		OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
		TokenKey:  tokenString[:16],
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(5) * time.Minute),
	}); result.Error != nil {
		H.logger(
			c, utils.AuthFirstFactor, result.Error.Error(), "", "error", "Failed create mfa token",
			user.Slug,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	var testOTP string

	if H.Conf.ENVIRONMENT == "testing" {
		testOTP = strings.Join(oneTimePasscode, "")
	} else if err = H.sendSMS(
		user.PhoneNumber, "One-time passcode:\n" + strings.Join(oneTimePasscode, " "),
	); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed send sms otp", user.Slug)

		if result := H.DBs.ApiGateway.Exec(
			"DELETE FROM mfa_tokens WHERE token_key = ?", tokenString[:16],
		); result.Error != nil {
			H.logger(
				c, utils.AuthFirstFactor, result.Error.Error(), "", "error", "Failed delete mfa token",
				user.Slug,
			)
		}

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).JSON(&AuthFirstFactorResponseBody{
		MFAToken: tokenString,
		TestOTP: 	testOTP,
	})
}

func (H Handler) deleteMfaToken(c *fiber.Ctx, token *models.MFAToken) {
	if result := H.DBs.ApiGateway.Delete(&token); result.Error != nil {
		H.logger(
			c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB,
			token.UserSlug,
		)
	} else if n := result.RowsAffected; n != 1 {
		H.logger(
			c, utils.AuthFirstFactor, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB, token.UserSlug,
		)
	}
}
