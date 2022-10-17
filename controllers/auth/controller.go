package auth

import (
	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

type handler struct {
	DBs  *databases.Databases
	Conf *config.AppConfig
}

func (h handler) createLog(
	c *fiber.Ctx,
	caller string,
	clientOperation string,
	detail string,
	extra string,
	level string,
	message string,
) {
	var clientIP string

	if h.Conf.GO_FIBER_BEHIND_PROXY {
		clientIP = c.Get("X-Forwarded-For")
	} else {
		clientIP = c.IP()
	}

	h.DBs.Logger.Create(&models.Log{
		Caller:          caller,
		ClientIP:        clientIP,
		ClientOperation: clientOperation,
		ContextString:   c.String(),
		Detail:          detail,
		Extra:           extra,
		Level:           level,
		Message:         message,
		RequestBody:     string(c.Body()),
	})
}

func RegisterAuth(api *fiber.Router, dbs *databases.Databases, conf *config.AppConfig) {
	h := handler{dbs, conf}
	authApi := (*api).Group("/auth")
	authApi.Post("/create_account", h.CreateAccount)
	authApi.Post("/log_in_account", h.LogInAccount)
}
