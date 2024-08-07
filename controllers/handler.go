package controllers

import (
	"runtime"
	"strconv"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/databases"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

type Handler struct {
	DBs  *databases.Databases
	Conf *config.AppConfig
}

func (H Handler) createLog(
	c *fiber.Ctx, caller, clientOperation, detail, extra, level, message string,
) {
	var clientIP string

	if H.Conf.BEHIND_PROXY {
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

func (H Handler) logger(c *fiber.Ctx, clientOperation, detail, extra, level, message string) {
	_, file, line, _ := runtime.Caller(1)

	H.createLog(
		c, file + ":" + strconv.FormatInt(int64(line), 10), clientOperation, detail, extra, level,
		message,
	)
}

func (H Handler) sendSMS(phoneNumber, messageBody string) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: H.Conf.TWILIO_ACCOUNT_SID,
		Password: H.Conf.TWILIO_AUTH_TOKEN,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(H.Conf.TWILIO_PHONE_NUMBER)
	params.SetTo(phoneNumber)
	params.SetBody(messageBody)

	if resp, err := client.Api.CreateMessage(params); err != nil {
		return err
	} else if response, err := json.Marshal(resp); err != nil {
		return err
	} else {
		if H.Conf.GO_TESTING_CONTEXT != nil {
			H.Conf.GO_TESTING_CONTEXT.Log("\n\nResponse:\n" + string(response) + "\n\n")
		}

		return nil
	}
}
