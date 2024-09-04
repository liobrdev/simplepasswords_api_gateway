package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type CreateSecretRequestBody struct {
	UserSlug			 string `json:"user_slug"`
	VaultSlug			 string `json:"vault_slug"`
	EntrySlug			 string `json:"entry_slug"`
	SecretLabel		 string `json:"secret_label"`
	SecretString	 string `json:"secret_string"`
	SecretPriority uint8	`json:"secret_priority"`
}

func (H Handler) VaultsCreateSecret(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.CreateSecret {
		H.logger(c, utils.CreateSecret, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := CreateSecretRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.CreateSecret, err.Error(), "", "error", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	agent := fiber.Post("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/secrets")
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.CreateSecret)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.CreateSecret, errString, "", "error", utils.ErrorVaultsCreateSecret)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
