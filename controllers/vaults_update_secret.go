package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type UpdateSecretRequestBody struct {
	Label  string `json:"secret_label"`
	String string `json:"secret_string"`
}

func (H Handler) VaultsUpdateSecret(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.UpdateSecret {
		H.logger(c, utils.UpdateSecret, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := UpdateSecretRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.UpdateSecret, err.Error(), "", "error", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	slug := c.Params("slug")

	agent := fiber.Patch(
		"http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/secrets/" + slug,
	)

	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.UpdateSecret)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.UpdateSecret, errString, "", "error", utils.ErrorVaultsUpdateSecret)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
