package controllers

import (
	"runtime"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

type Handler struct {
	DBs  *databases.Databases
	Conf *config.AppConfig
}

func (H Handler) createLog(c *fiber.Ctx, caller string, clientOperation string, detail string,
extra string, level string, message string) {
	var clientIP string

	if H.Conf.GO_FIBER_BEHIND_PROXY {
		clientIP = c.Get("X-Forwarded-For")
	} else {
		clientIP = c.IP()
	}

	H.DBs.Logger.Create(&models.Log{
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

func (H Handler) logger(c *fiber.Ctx, clientOperation string, detail string, extra string,
level string, message string) {
	_, file, line, _ := runtime.Caller(1)

	H.createLog(
		c, file + ":" + strconv.FormatInt(int64(line), 10), clientOperation, detail, extra, level,
		message,
	)
}
