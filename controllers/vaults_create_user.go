package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) VaultsCreateUser(userSlug string) string {
	agent := fiber.Post("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/users")
	agent.Set("Content-Type", "application/json")
	agent.Set("Client-Operation", utils.CreateUser)
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)

	agent.JSON(fiber.Map{ "user_slug": userSlug })
	statusCode, body, errs := agent.String()

	var errorString string

	if len(errs) > 0 {
		for _, err := range errs {
			errorString += err.Error() + ";;"
		}

		if body != "" {
			errorString += body + ";;"
		}
	} else if statusCode != 204 && body != "" {
		errorString += body + ";;"
	}

	return errorString
}
