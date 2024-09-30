package controllers

import (
	"bytes"
	"net/smtp"
	"runtime"
	"strconv"
	"strings"
	"text/template"

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

	if _, err := client.Api.CreateMessage(params); err != nil {
		return err
	}

	return nil
}

func (H Handler) sendEmail(
	subject, from string, to []string, templateFile string, data map[string]string,
) error {
	auth := smtp.PlainAuth("", H.Conf.AWS_SES_KEY, H.Conf.AWS_SES_PASSWORD, H.Conf.EMAIL_HOST)

	var t *template.Template
	var err error

	if t, err = template.ParseFiles(templateFile); err != nil {
		return err
	}

	var body bytes.Buffer

	if _, err = body.Write([]byte(
		"From: " + from + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=utf-8" + "\r\n\r\n",
	)); err != nil {
		return err
	}

	t.Execute(&body, data)

	if _, err = body.Write([]byte("\r\n")); err != nil {
		return err
	}

	return smtp.SendMail(H.Conf.EMAIL_HOST + ":" + H.Conf.EMAIL_PORT, auth, from, to, body.Bytes())
}
