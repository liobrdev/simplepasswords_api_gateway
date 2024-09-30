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
		H.logger(c, utils.AuthFirstFactor, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	body := AuthFirstFactorRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.EmailRegexp.Match([]byte(body.Email)) {
		H.logger(c, utils.AuthFirstFactor, body.Email, "", "warn", utils.ErrorAcctEmail)

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.AuthFirstFactor, "", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	var user models.User

	if result := H.DBs.ApiGateway.Where("email_address = ?", body.Email).Limit(1).Find(&user);
	result.Error != nil {
		H.logger(c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	} else if n != 1 {
		H.logger(
			c, utils.AuthFirstFactor, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if !user.IsActive {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	if password, err := hex.DecodeString(body.Password); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	} else if !utils.CompareHashAndPassword(user.PasswordHash, password, user.PasswordSalt) {
		return utils.RespondWithError(c, 400, utils.ErrorFailedLogin, nil, nil)
	}

	var userAllMFATokens []models.MFAToken

	// Get all user's MFA tokens for cleanup
	if result := H.DBs.ApiGateway.Where("user_slug = ?", user.Slug).Find(&userAllMFATokens);
	result.Error != nil {
		H.logger(c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB)
	}

	now := time.Now().UTC()

	// Clean up expired MFA tokens
	for _, token := range userAllMFATokens {
		if !token.ExpiresAt.After(now) {
			if result := H.DBs.ApiGateway.Delete(&token); result.Error != nil {
				H.logger(
					c, utils.AuthFirstFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB,
				)
			} else if n := result.RowsAffected; n != 1 {
				H.logger(
					c, utils.AuthFirstFactor, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
					"error", utils.ErrorFailedDB,
				)
			}
		}
	}

	var mfaTokenString	string
	var oneTimePasscode []string
	var err error

	if mfaTokenString, err = utils.GenerateSlug(80); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed generate mfa string")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if oneTimePasscode, err = utils.GenerateOTP(); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed generate otp")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if result := H.DBs.ApiGateway.Create(&models.MFAToken{
		UserSlug:  user.Slug,
		KeyDigest: utils.HashToken(mfaTokenString),
		OTPDigest: utils.HashToken(strings.Join(oneTimePasscode, "")),
		TokenKey:  mfaTokenString[:16],
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(5) * time.Minute),
	}); result.Error != nil {
		H.logger(
			c, utils.AuthFirstFactor, result.Error.Error(), "", "error", "Failed create mfa token",
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	var testOTP string

	if H.Conf.ENVIRONMENT == "testing" {
		testOTP = strings.Join(oneTimePasscode, "")
	} else if err = H.sendSMS(
		user.PhoneNumber, "One-time passcode:\n" + strings.Join(oneTimePasscode, " "),
	); err != nil {
		H.logger(c, utils.AuthFirstFactor, err.Error(), "", "error", "Failed send sms otp")

		if result := H.DBs.ApiGateway.Exec(
			"DELETE FROM mfa_tokens WHERE token_key = ?", mfaTokenString[:16],
		); result.Error != nil {
			H.logger(
				c, utils.AuthFirstFactor, result.Error.Error(), "", "error", "Failed delete mfa token",
			)
		}

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).JSON(&AuthFirstFactorResponseBody{
		MFAToken: mfaTokenString,
		TestOTP: 	testOTP,
	})
}
