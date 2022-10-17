package auth

import (
	"errors"
	// "io"
	// "net/http"
	"strconv"
	// "strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
	// "github.com/liobrdev/simplepasswords_api_gateway/validators"
)

type LogInAccountRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInAccountResponseBody struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

const LogInAccountCaller string = "api_gateway/controllers/auth/log_in_account.go"

func (h handler) LogInAccount(c *fiber.Ctx) error {
	body := LogInAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, err.Error(), "", "warn",
			utils.ErrorParse,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email == "" || len(body.Email) > 256 {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, body.Email, "", "warn",
			utils.ErrorAcctEmail,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, "", "", "warn", utils.ErrorAcctPW,
		)

		return utils.RespondWithError(c, fiber.StatusBadRequest, utils.ErrorBadRequest, nil, nil)
	}

	var user models.User

	if result := h.DBs.ApiGateway.First(
		&user, "email_address = ?", body.Email,
	); result.Error != nil {
		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			h.createLog(
				c, LogInAccountCaller, utils.LogInAccount, result.Error.Error(), "", "error",
				utils.ErrorFailedDB,
			)
		}

		return utils.RespondWithError(c, fiber.StatusNotFound, utils.ErrorFailedLogin, nil, nil)
	} else if n := result.RowsAffected; n != 1 {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, "result.RowsAffected != 1",
			strconv.FormatInt(n, 10), "error", utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, fiber.StatusNotFound, utils.ErrorFailedLogin, nil, nil)
	}

	if err := bcrypt.CompareHashAndPassword(
		user.PasswordHash, []byte(body.Password+user.PasswordSalt),
	); err != nil {
		return utils.RespondWithError(c, fiber.StatusNotFound, utils.ErrorFailedLogin, nil, nil)
	}

	var sessionToken string
	var err error

	if sessionToken, err = utils.GenerateSlug(144); err != nil {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, err.Error(), "", "error",
			"Failed generate session token",
		)

		return utils.RespondWithError(
			c, fiber.StatusInternalServerError, utils.ErrorServer, nil, nil,
		)
	}

	createdAt := time.Now().UTC()

	if result := h.DBs.ApiGateway.Create(&models.ClientSession{
		UserSlug:  user.Slug,
		Digest:    utils.HashToken(sessionToken),
		TokenKey:  sessionToken[:16],
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(time.Duration(15) * time.Minute),
	}); result.Error != nil {
		h.createLog(
			c, LogInAccountCaller, utils.LogInAccount, result.Error.Error(), "", "error",
			"Failed create client session",
		)

		return utils.RespondWithError(
			c, fiber.StatusInternalServerError, utils.ErrorServer, nil, nil,
		)
	}

	return c.Status(fiber.StatusOK).JSON(&LogInAccountResponseBody{
		Token: sessionToken,
		User:  user,
	})
}
