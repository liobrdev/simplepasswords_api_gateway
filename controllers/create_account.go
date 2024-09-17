package controllers

import (
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/config"
	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type CreateAccountRequestBody struct {
	Name      	string `json:"name"`
	Email     	string `json:"email"`
	Phone				string `json:"phone"`
	Password		string `json:"password"`
}

type CreateAccountResponseBody struct {
	Token	string			`json:"token"`
	User	models.User	`json:"user"`
}

func (H Handler) CreateAccount(c *fiber.Ctx) error {
	if header := c.Get("Client-Operation"); header != utils.CreateAccount {
		H.logger(c, utils.CreateAccount, header, "", "warn", utils.ErrorClientOperation)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	body := CreateAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.NameRegexp.Match([]byte(body.Name)) {
		H.logger(c, utils.CreateAccount, body.Name, "", "warn", utils.ErrorAcctName)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.EmailRegexp.Match([]byte(body.Email)) {
		H.logger(c, utils.CreateAccount, body.Email, "", "warn", utils.ErrorAcctEmail)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.PhoneRegexp.Match([]byte(body.Phone)) {
		H.logger(c, utils.CreateAccount, body.Phone, "", "warn", utils.ErrorAcctPhone)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.CreateAccount, "", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var user models.User

	if password, err := hex.DecodeString(body.Password); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed decode password")

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	} else if hash, salt, err := utils.GenerateUserCredentials(password); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed generate user credentials")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else {
		user.PasswordHash = hash
		user.PasswordSalt = salt
	}

	if userSlug, err := utils.GenerateSlug(32); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed generate user.Slug")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else {
		user.Slug = userSlug
	}

	user.Name = body.Name
	user.EmailAddress = body.Email
	user.PhoneNumber = body.Phone

	var sessionToken string

	if err := H.DBs.ApiGateway.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&user); result.Error != nil {
			return result.Error
		}

		var err error

		if sessionToken, err = utils.GenerateSlug(80); err != nil {
			return err
		}

		createdAt := time.Now().UTC()

		if result := tx.Create(&models.ClientSession{
			UserSlug:  user.Slug,
			ClientIP:  c.IP(),
			Digest:    utils.HashToken(sessionToken),
			TokenKey:  sessionToken[:16],
			CreatedAt: createdAt,
			ExpiresAt: createdAt.Add(time.Duration(15) * time.Minute),
		}); result.Error != nil {
			return result.Error
		}

		return nil
	}); err != nil {
		if !utils.UniqueConstraintRegexp.Match([]byte(err.Error())) {
			H.logger(c, utils.CreateAccount, err.Error(), "", "error", utils.ErrorCreateUser)
		}

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if H.Conf.ENVIRONMENT != "testing" {
		if errorString := vaultsCreateUser(H.Conf, user.Slug); errorString != "" {
			H.logger(c, utils.CreateAccount, errorString, "", "error", utils.ErrorVaultsCreateUser)

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(&CreateAccountResponseBody{
		Token: sessionToken,
		User:	 models.User{ Slug: user.Slug, Name: user.Name },
	})
}

func vaultsCreateUser(conf *config.AppConfig, userSlug string) string {
	agent := fiber.Post("http://" + conf.VAULTS_HOST + ":" + conf.VAULTS_PORT + "/api/users")
	agent.Set("Content-Type", "application/json")
	agent.Set("Client-Operation", utils.CreateUser)
	agent.Set("Authorization", "Token " + conf.VAULTS_ACCESS_TOKEN)

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
