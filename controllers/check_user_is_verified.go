package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) CheckUserIsVerified(c *fiber.Ctx) error {
	var session *models.ClientSession
	var ok bool

	if session, ok = c.UserContext().Value(sessionContextKey{}).(*models.ClientSession); !ok {
		H.logger(c, c.Get("Client-Operation"), "", "", "error", "Failed session.User context", "")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if !session.User.EmailIsVerified {
		H.logger(
			c, c.Get("Client-Operation"), "", "", "error", "User email not verified", session.UserSlug,
		)

		return utils.RespondWithError(c, 403, utils.ErrorServer, nil, nil)
	}

	if !session.User.PhoneIsVerified {
		H.logger(
			c, c.Get("Client-Operation"), "", "", "error", "User phone not verified", session.UserSlug,
		)

		return utils.RespondWithError(c, 403, utils.ErrorServer, nil, nil)
	}

	return c.Next()
}
