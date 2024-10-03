package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type UpdateVaultRequestBody struct {
	Title string `json:"vault_title"`
}

func (H Handler) VaultsUpdateVault(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.UpdateVault {
		H.logger(c, utils.UpdateVault, header, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := UpdateVaultRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.UpdateVault, err.Error(), "", "error", utils.ErrorParse, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Patch(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/vaults/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.UpdateVault)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.UpdateVault, errString, "", "error", utils.ErrorVaultsUpdateVault, "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
