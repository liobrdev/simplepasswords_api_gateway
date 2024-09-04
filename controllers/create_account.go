package controllers

import (
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
	Password  	string `json:"password"`
	Password2 	string `json:"password_2"`
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

	if body.Name != H.Conf.ADMIN_NAME {
		if body.Name == "" {
			H.logger(c, utils.CreateAccount, body.Name, "", "warn", utils.ErrorAcctName)
		}

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email != H.Conf.ADMIN_EMAIL {
		if body.Email == "" {
			H.logger(c, utils.CreateAccount, body.Email, "", "warn", utils.ErrorAcctEmail)
		}

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Phone != H.Conf.ADMIN_PHONE {
		if body.Phone == "" {
			H.logger(c, utils.CreateAccount, body.Phone, "", "warn", utils.ErrorAcctPhone)
		}

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.CreateAccount, "", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password != body.Password2 {
		H.logger(c, utils.CreateAccount, "", "", "warn", utils.ErrorNonMatchPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == body.Email {
		H.logger(c, utils.CreateAccount, "Password is email", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.Password) > 72 - len(H.Conf.ADMIN_SALT_1) - len(H.Conf.ADMIN_SALT_2) {
		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.Password) < 16 {
		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.ContainsUppercase(body.Password) {
		H.logger(c, utils.CreateAccount, "Missing uppercase", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.ContainsLowercase(body.Password) {
		H.logger(c, utils.CreateAccount, "Missing lowercase", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.ContainsNumber(body.Password) {
		H.logger(c, utils.CreateAccount, "Missing number", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if !utils.ContainsSpecialChar(body.Password) {
		H.logger(c, utils.CreateAccount, "Missing special char", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if utils.ContainsWhitespace(body.Password) {
		H.logger(c, utils.CreateAccount, "Has whitespace", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var user models.User

	if userSlug, err := utils.GenerateSlug(32); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed generate user.Slug")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else {
		user.Slug = userSlug
	}

	if hash, err := utils.GenerateUserCredentials(body.Password, H.Conf); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed generate user credentials")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else {
		user.PasswordHash = hash
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
