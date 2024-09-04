package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsDeleteEntry(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.DeleteEntry {
		H.logger(c, utils.DeleteEntry, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Delete(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/entries/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.DeleteEntry)
	agent.Set("Content-Type", "application/json")

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.DeleteEntry, errString, "", "error", utils.ErrorVaultsDeleteEntry)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
