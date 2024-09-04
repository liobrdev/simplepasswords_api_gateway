package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type CreateVaultRequestBody struct {
	UserSlug   string `json:"user_slug"`
	VaultTitle string `json:"vault_title"`
}

func (H Handler) VaultsCreateVault(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.CreateVault {
		H.logger(c, utils.CreateVault, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := CreateVaultRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.CreateVault, err.Error(), "", "error", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	agent := fiber.Post("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/vaults")
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.CreateVault)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.CreateVault, errString, "", "error", utils.ErrorVaultsCreateVault)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
