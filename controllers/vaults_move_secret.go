package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type MoveSecretRequestBody struct {
	Priority 	string `json:"secret_priority"`
	EntrySlug string `json:"entry_slug"`
}

func (H Handler) VaultsMoveSecret(c *fiber.Ctx) error {
	if clientOperation := c.Get("Client-Operation"); clientOperation != utils.MoveSecret {
		H.logger(c, utils.MoveSecret, clientOperation, "", "warn", utils.ErrorClientOperation, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := MoveSecretRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.MoveSecret, err.Error(), "", "error", utils.ErrorParse, "")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Patch(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/secrets/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.MoveSecret)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.MoveSecret, errString, "", "error", utils.ErrorVaultsMoveSecret, "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
