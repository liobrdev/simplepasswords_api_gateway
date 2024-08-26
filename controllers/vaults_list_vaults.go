package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsListVaults(c *fiber.Ctx) error {
	var user *models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(*models.User); !ok {
		H.logger(c, utils.RetrieveUser, "", "", "error", "Failed user context")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	agent := fiber.Get("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/vaults")
	agent.Set("Content-Type", "application/json")
	agent.Set("Client-Operation", utils.ListVaults)
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("User-Slug", user.Slug )

	statusCode, body, errs := agent.String()

	var errorString string

	if len(errs) > 0 {
		for _, err := range errs {
			errorString += err.Error() + ";;"
		}

		if body != "" {
			errorString += body + ";;"
		}
	}

	if statusCode != 200 {
		if body != "" {
			errorString += body + ";;"
		}

		H.logger(c, utils.ListVaults, errorString, "", "error", utils.ErrorVaultsListVaults)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).Send([]byte(body))
}