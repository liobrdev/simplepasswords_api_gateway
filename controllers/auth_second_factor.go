package controllers

import (
	"bytes"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type AuthSecondFactorRequestBody struct {
	MFAToken string `json:"mfa_token"`
	PhoneOTP string `json:"phone_otp"`
}

type AuthSecondFactorResponseBody struct {
	Token string			`json:"token"`
	User  models.User	`json:"user"`
}

func (H Handler) AuthSecondFactor(c *fiber.Ctx) error {
	body := AuthSecondFactorRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.AuthSecondFactor, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.MFAToken) != 80 {
		H.logger(c, utils.AuthSecondFactor, body.MFAToken, "", "warn", utils.ErrorMFAToken)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.PhoneOTP) != 20 {
		H.logger(c, utils.AuthSecondFactor, body.PhoneOTP, "", "warn", utils.ErrorPhoneOTP)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var mfaToken models.MFAToken
	keyDigest := utils.HashToken(body.MFAToken)

	if result := H.DBs.ApiGateway.Preload("User").Where("key_digest = ?", keyDigest).Limit(1).
	Find(&mfaToken); result.Error != nil {
		H.logger(c, utils.AuthSecondFactor, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	} else if n != 1 {
		H.logger(
			c, utils.AuthSecondFactor, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	}

	now := time.Now().UTC()

	if !mfaToken.ExpiresAt.After(now) {
		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	}

	if !mfaToken.User.IsActive {
		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	}

	if !bytes.Equal(mfaToken.OTPDigest, utils.HashToken(body.PhoneOTP)) {
		return utils.RespondWithError(c, 400, utils.ErrorAuthenticate, nil, nil)
	}

	var sessionToken string
	var err error

	if sessionToken, err = utils.GenerateSlug(80); err != nil {
		H.logger(c, utils.AuthSecondFactor, err.Error(), "", "error", "Failed generate session token")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if result := H.DBs.ApiGateway.Create(&models.ClientSession{
		UserSlug:  mfaToken.UserSlug,
		ClientIP:  c.IP(),
		Digest:    utils.HashToken(sessionToken),
		TokenKey:  sessionToken[:16],
		CreatedAt: now,
		ExpiresAt: now.Add(time.Duration(15) * time.Minute),
	}); result.Error != nil {
		H.logger(
			c, utils.AuthSecondFactor, result.Error.Error(), "", "error", "Failed create client session",
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	hiddenUser := models.User{
		Slug: 				mfaToken.UserSlug,
		Name: 				mfaToken.User.Name,
		EmailAddress:	utils.HideEmail(mfaToken.User.EmailAddress),
		PhoneNumber:	utils.HidePhone(mfaToken.User.PhoneNumber),
		IsActive: 		mfaToken.User.IsActive,
	}

	return c.Status(fiber.StatusOK).JSON(&AuthSecondFactorResponseBody{
		Token: sessionToken,
		User:  hiddenUser,
	})
}
