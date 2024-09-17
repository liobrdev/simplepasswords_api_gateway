package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

func (H Handler) LogoutAccount(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.LogoutAccount {
		H.logger(c, utils.LogoutAccount, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var session *models.ClientSession
	var ok bool

	if session, ok = c.UserContext().Value(sessionContextKey{}).(*models.ClientSession); !ok {
		H.logger(c, utils.LogoutAccount, "", "", "error", "Failed session.User context")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if result := H.DBs.ApiGateway.Delete(&session); result.Error != nil {
		H.logger(c, utils.LogoutAccount, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if n := result.RowsAffected; n != 1 {
		H.logger(
			c, utils.LogoutAccount, "result.RowsAffected != 1", strconv.FormatInt(n, 10),
			"error", utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
