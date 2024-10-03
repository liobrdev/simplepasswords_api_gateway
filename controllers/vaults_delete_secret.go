package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsDeleteSecret(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.DeleteSecret {
		H.logger(c, utils.DeleteSecret, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Delete(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/secrets/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.DeleteSecret)
	agent.Set("Content-Type", "application/json")

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.DeleteSecret, errString, "", "error", utils.ErrorVaultsDeleteSecret, "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
