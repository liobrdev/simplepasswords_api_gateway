package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsListVaults(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.ListVaults {
		H.logger(c, utils.ListVaults, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var user *models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(*models.User); !ok {
		H.logger(c, utils.ListVaults, "", "", "error", "Failed user context")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	agent := fiber.Get("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/vaults")
	agent.Set("Content-Type", "application/json")
	agent.Set("Client-Operation", utils.ListVaults)
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("User-Slug", user.Slug )

	_, body, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.ListVaults, errString, "", "error", utils.ErrorVaultsListVaults)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).Send([]byte(body))
}
