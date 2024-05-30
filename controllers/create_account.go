package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type CreateAccountRequestBody struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
}

type CreateAccountResponseBody struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (H Handler) CreateAccount(c *fiber.Ctx) error {
	body := CreateAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Name == "" || len(body.Name) > 64 {
		H.logger(c, utils.CreateAccount, body.Name, "", "warn", utils.ErrorAcctName)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email == "" || len(body.Email) > 256 {
		H.logger(c, utils.CreateAccount, body.Email, "", "warn", utils.ErrorAcctEmail)

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

	if len(body.Password) > 256 {
		H.logger(
			c, utils.CreateAccount, "Too long: " + strconv.Itoa(len(body.Password)) + " > 256", "",
			"warn", utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.Password) < 16 {
		H.logger(
			c, utils.CreateAccount, "Too short: " + strconv.Itoa(len(body.Password)) + " < 16", "",
			"warn", utils.ErrorAcctPW,
		)

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

	if salt, hash, err := utils.GenerateUserCredentials(body.Password); err != nil {
		H.logger(c, utils.CreateAccount, err.Error(), "", "error", "Failed generate user credentials")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else {
		user.PasswordSalt = salt
		user.PasswordHash = hash
	}

	user.Name = body.Name
	user.EmailAddress = body.Email

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
		if err.Error() == "UNIQUE constraint failed: users.email_address" {
			return utils.RespondWithError(c, fiber.StatusConflict, utils.ErrorDiffEmail, nil, nil)
		} else {
			H.logger(c, utils.CreateAccount, err.Error(), "", "error", utils.ErrorCreateUser)

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		}
	}

	if H.Conf.GO_FIBER_ENVIRONMENT != "testing" {
		if err := H.vaultsCreateUser(user.Slug); err != nil {
			H.logger(c, utils.CreateAccount, err.Error(), "", "error", utils.ErrorVaultsCreateUser)

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(&CreateAccountResponseBody{
		Token: sessionToken,
		User:  user,
	})
}

func (H Handler) vaultsCreateUser(userSlug string) error {
	reqBody := strings.NewReader(`{"user_slug":"` + userSlug + `"}`)

	if resp, err := http.Post(H.Conf.GO_FIBER_VAULTS_URL, "application/json", reqBody); err != nil {
		return err
	} else if resp.StatusCode != fiber.StatusNoContent {
		if respBody, err := io.ReadAll(resp.Body); err != nil {
			return err
		} else {
			return errors.New(string(respBody))
		}
	}

	return nil
}
