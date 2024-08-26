package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type RetrieveUserResponseBody struct {
	Name string `json:"user_name"`
}

func (H Handler) RetrieveUser(c *fiber.Ctx) error {
	var user *models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(*models.User); !ok {
		H.logger(c, utils.RetrieveUser, "", "", "error", "Failed user context")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).JSON(&RetrieveUserResponseBody{ Name: user.Name })
}
