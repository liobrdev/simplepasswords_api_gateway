package controllers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type LogInAccountRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogInAccountResponseBody struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (H Handler) LogInAccount(c *fiber.Ctx) error {
	body := LogInAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.LogInAccount, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email == "" || len(body.Email) > 256 {
		H.logger(c, utils.LogInAccount, body.Email, "", "warn",
			utils.ErrorAcctEmail,
		)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.LogInAccount, "", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	var user models.User

	if result := H.DBs.ApiGateway.Where("email_address = ?", body.Email).Limit(1).Find(&user);
	result.Error != nil {
		H.logger(c, utils.LogInAccount, result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 404, utils.ErrorFailedLogin, nil, nil)
	} else if n := result.RowsAffected; n == 0 {
		return utils.RespondWithError(c, 404, utils.ErrorFailedLogin, nil, nil)
	} else if n != 1 {
		H.logger(
			c, utils.LogInAccount, "result.RowsAffected != 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 404, utils.ErrorFailedLogin, nil, nil)
	}

	if err := bcrypt.CompareHashAndPassword(
		user.PasswordHash, []byte(body.Password+user.PasswordSalt),
	); err != nil {
		return utils.RespondWithError(c, 404, utils.ErrorFailedLogin, nil, nil)
	}

	var sessionToken string
	var err error

	if sessionToken, err = utils.GenerateSlug(80); err != nil {
		H.logger(c, utils.LogInAccount, err.Error(), "", "error", "Failed generate session token")

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	createdAt := time.Now().UTC()

	if result := H.DBs.ApiGateway.Create(&models.ClientSession{
		UserSlug:  user.Slug,
		ClientIP:  c.IP(),
		Digest:    utils.HashToken(sessionToken),
		TokenKey:  sessionToken[:16],
		CreatedAt: createdAt,
		ExpiresAt: createdAt.Add(time.Duration(15) * time.Minute),
	}); result.Error != nil {
		H.logger(
			c, utils.LogInAccount, result.Error.Error(), "", "error", "Failed create client session",
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	return c.Status(fiber.StatusOK).JSON(&LogInAccountResponseBody{
		Token: sessionToken,
		User:  user,
	})
}
