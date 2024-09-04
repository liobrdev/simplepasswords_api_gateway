package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type reqBodySecret struct {
	Label  	 string	`json:"secret_label"`
	String 	 string	`json:"secret_string"`
	Priority uint8 	`json:"secret_priority"`
}

type CreateEntryRequestBody struct {
	UserSlug   string          `json:"user_slug"`
	VaultSlug  string          `json:"vault_slug"`
	EntryTitle string          `json:"entry_title"`
	Secrets    []reqBodySecret `json:"secrets"`
}

func (H Handler) VaultsCreateEntry(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.CreateEntry {
		H.logger(c, utils.CreateEntry, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	reqBody := CreateEntryRequestBody{}

	if err := c.BodyParser(&reqBody); err != nil {
		H.logger(c, utils.CreateEntry, err.Error(), "", "error", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	agent := fiber.Post("http://" + H.Conf.VAULTS_HOST + ":" + H.Conf.VAULTS_PORT + "/api/entries")
	agent.Set("Authorization", "Token " + H.Conf.VAULTS_ACCESS_TOKEN)
	agent.Set("Client-Operation", utils.CreateEntry)
	agent.JSON(&reqBody)

	_, _, errString := checkVaultsResponse(agent)

	if errString != "" {
		H.logger(c, utils.CreateEntry, errString, "", "error", utils.ErrorVaultsCreateEntry)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.SendStatus(204)
}
