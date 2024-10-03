package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VerifyPhoneConfirm(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.VerifyPhoneConfirm {
		H.logger(c, utils.VerifyPhoneConfirm, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	return c.SendStatus(204)
}
