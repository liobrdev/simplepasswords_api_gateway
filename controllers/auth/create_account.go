package auth

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

const CreateAccountCaller string = "api_gateway/controllers/auth/create_account.go"

func (h handler) CreateAccount(c *fiber.Ctx) error {
	body := CreateAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, err.Error(), "", "warn",
			utils.ErrorParse,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Name == "" || len(body.Name) > 64 {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, body.Name, "", "warn",
			utils.ErrorAcctName,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email == "" || len(body.Email) > 256 {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, body.Email, "", "warn",
			utils.ErrorAcctEmail,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "", "", "warn", utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password != body.Password2 {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "", "", "warn", utils.ErrorNonMatchPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == body.Email {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Password is email", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.Password) > 256 {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount,
			"Too long: "+strconv.Itoa(len(body.Password))+" > 256", "", "warn", utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if len(body.Password) < 16 {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount,
			"Too short: "+strconv.Itoa(len(body.Password))+" < 16", "", "warn", utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if !ContainsUppercase(body.Password) {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Missing uppercase", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if !ContainsLowercase(body.Password) {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Missing lowercase", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if !ContainsNumber(body.Password) {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Missing number", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if !ContainsSpecialChar(body.Password) {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Missing special char", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if ContainsWhitespace(body.Password) {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, "Has whitespace", "", "warn",
			utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	var user models.User

	if userSlug, err := utils.GenerateSlug(32); err != nil {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, err.Error(), "", "error",
			"Failed generate user.Slug",
		)

		return utils.RespondWithError(
			c, fiber.StatusInternalServerError, utils.ErrorServer, nil, nil,
		)
	} else {
		user.Slug = userSlug
	}

	if salt, hash, err := utils.GenerateUserCredentials(body.Password); err != nil {
		h.createLog(
			c, CreateAccountCaller, utils.CreateAccount, err.Error(), "", "error",
			"Failed generate user credentials",
		)

		return utils.RespondWithError(
			c, fiber.StatusInternalServerError, utils.ErrorServer, nil, nil,
		)
	} else {
		user.PasswordSalt = salt
		user.PasswordHash = hash
	}

	user.Name = body.Name
	user.EmailAddress = body.Email

	var sessionToken string

	if err := h.DBs.ApiGateway.Transaction(func(tx *gorm.DB) error {
		if result := tx.Create(&user); result.Error != nil {
			return result.Error
		}

		var err error

		if h.Conf.GO_FIBER_ENVIRONMENT != "testing" {
			if err = h.vaultsCreateUser(user.Slug); err != nil {
				return err
			}
		}

		if sessionToken, err = utils.GenerateSlug(144); err != nil {
			return err
		}

		createdAt := time.Now().UTC()

		if result := tx.Create(&models.ClientSession{
			UserSlug:  user.Slug,
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
			h.createLog(
				c, CreateAccountCaller, utils.CreateAccount, err.Error(), "", "error",
				utils.ErrorCreateUser,
			)

			return utils.RespondWithError(
				c, fiber.StatusInternalServerError, utils.ErrorServer, nil, nil,
			)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(&CreateAccountResponseBody{
		Token: sessionToken,
		User:  user,
	})
}

func (h handler) vaultsCreateUser(userSlug string) error {
	reqBody := strings.NewReader(`{"user_slug":"` + userSlug + `"}`)

	if resp, err := http.Post(
		h.Conf.GO_FIBER_VAULTS_URL, "application/json", reqBody,
	); err != nil {
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
