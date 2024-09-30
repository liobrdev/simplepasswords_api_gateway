package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) RetrieveUser(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.RetrieveUser {
		H.logger(c, utils.RetrieveUser, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var session *models.ClientSession
	var ok bool

	if session, ok = c.UserContext().Value(sessionContextKey{}).(*models.ClientSession); !ok {
		H.logger(c, utils.RetrieveUser, "", "", "error", "Failed session.User context")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(200).JSON(&models.User{
		Slug: session.UserSlug,
		Name: session.User.Name,
		EmailIsVerified: session.User.EmailIsVerified,
		PhoneIsVerified: session.User.PhoneIsVerified,
	})
}
