package controllers

import (
	"bytes"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type VerifyEmailConfirmRequestBody struct {
	Token string `json:"token"`
}

func (H Handler) VerifyEmailConfirm(c *fiber.Ctx) error {
	var user models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(models.User); !ok {
		H.logger(c, utils.VerifyEmailConfirm, "UserContext not ok", "", "error", utils.ErrorUserContext)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if user.EmailIsVerified {
		return c.SendStatus(fiber.StatusNoContent)
	}

	body := VerifyEmailConfirmRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.VerifyEmailConfirm, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.TokenRegexp.Match([]byte(body.Token)) {
		H.logger(c, utils.VerifyEmailConfirm, "token", body.Token, "warn", utils.ErrorToken)

		return utils.RespondWithError(c, 400, utils.ErrorToken, nil, nil)
	}

	var token models.EmailVerificationToken

	if result := H.DBs.ApiGateway.
	Where("user_slug = ? AND token_key = ?", user.Slug, body.Token[:16]).Limit(1).Find(&token);
	result.Error != nil {
		H.logger(c, utils.VerifyEmailConfirm, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if result.RowsAffected == 0 {
		return utils.RespondWithError(c, 400, utils.ErrorToken, nil, nil)
	}

	now := time.Now().UTC()

	if token.ExpiresAt.After(now) && bytes.Equal(utils.HashToken(body.Token), token.Digest) {
		user.EmailIsVerified = true
		H.DBs.ApiGateway.Save(&user)

		if result := H.DBs.ApiGateway.Where("user_slug = ?", user.Slug).
		Delete(&models.EmailVerificationToken{}); result.Error != nil {
			H.logger(c, utils.VerifyEmailConfirm, result.Error.Error(), "", "error", utils.ErrorFailedDB)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}

	return utils.RespondWithError(c, 400, utils.ErrorToken, nil, nil)
}
