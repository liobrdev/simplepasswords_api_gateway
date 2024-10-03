package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsRetrieveVault(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.RetrieveVault {
		H.logger(c, utils.RetrieveVault, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Get(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/vaults/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.RetrieveVault)
	agent.Set("Content-Type", "application/json")

	_, body, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.RetrieveVault, errString, "", "error", utils.ErrorVaultsRetrieveVault, "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).Send([]byte(body))
}
