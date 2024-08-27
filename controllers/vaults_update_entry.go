package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type UpdateEntryRequestBody struct {
	Title string `json:"entry_title"`
}

func (H Handler) VaultsUpdateEntry(c *fiber.Ctx) error {
	reqBody := UpdateEntryRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.UpdateEntry, err.Error(), "", "error", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Patch(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/entries/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.UpdateEntry)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.UpdateEntry, errString, "", "error", utils.ErrorVaultsUpdateEntry)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
