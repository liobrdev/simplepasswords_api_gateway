package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsRetrieveEntry(c *fiber.Ctx) error {
	slug := c.Params("slug")

	agent := fiber.Get(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/entries/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.RetrieveEntry)
	agent.Set("Content-Type", "application/json")

	_, body, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.RetrieveEntry, errString, "", "error", utils.ErrorVaultsRetrieveEntry)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).Send([]byte(body))
}