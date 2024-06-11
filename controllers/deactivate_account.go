package controllers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type DeactivateAccountRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (H Handler) DeactivateAccount(c *fiber.Ctx) error {
	urlParamSlug := c.Params("slug")

	if urlParamSlug == "" {
		H.logger(c, utils.DeactivateAccount, "urlParamSlug empty", "", "error", utils.ErrorParams)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	var user models.User
	var ok bool

	if user, ok = c.UserContext().Value(userContextKey{}).(models.User); !ok {
		H.logger(c, utils.DeactivateAccount, "UserContext not ok", "", "error", utils.ErrorUserContext)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}
	
	if urlParamSlug != user.Slug {
		H.logger(
			c, utils.DeactivateAccount, "urlParamSlug != user.Slug", urlParamSlug + ":" + user.Slug,
			"error", utils.ErrorParams,
		)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	body := DeactivateAccountRequestBody{}

	if err := c.BodyParser(&body); err != nil {
		H.logger(c, utils.DeactivateAccount, err.Error(), "", "warn", utils.ErrorParse)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email == "" || len(body.Email) > 256 {
		H.logger(c, utils.DeactivateAccount, body.Email, "", "warn", utils.ErrorAcctEmail)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Password == "" {
		H.logger(c, utils.DeactivateAccount, "", "", "warn", utils.ErrorAcctPW)

		return utils.RespondWithError(c, 400, utils.ErrorBadRequest, nil, nil)
	}

	if body.Email != user.EmailAddress {
		return utils.RespondWithError(c, 400, utils.ErrorFailedDeactivate, nil, nil)
	}

	if err := bcrypt.CompareHashAndPassword(
		user.PasswordHash, []byte(body.Password+user.PasswordSalt),
	); err != nil {
		return utils.RespondWithError(c, 400, utils.ErrorFailedDeactivate, nil, nil)
	}

	if err := H.DBs.ApiGateway.Transaction(func(tx *gorm.DB) error {
		if result := tx.Exec("DELETE FROM client_sessions WHERE user_slug = ?", user.Slug);
		result.Error != nil {
			return result.Error
		}

		if result := tx.Delete(&user); result.Error != nil {
			return result.Error
		} else if n := result.RowsAffected; n != 1 {
			return errors.New("deactivate_user; result.RowsAffected != 1; " + strconv.FormatInt(n, 10))
		}

		return nil
	}); err != nil {
		H.logger(c, utils.DeactivateAccount, err.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	if H.Conf.GO_FIBER_ENVIRONMENT != "testing" {
		if err := H.vaultsDeleteUser(user.Slug); err != nil {
			H.logger(c, utils.DeactivateAccount, err.Error(), "", "error", utils.ErrorVaultsDeleteUser)

			return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (H Handler) vaultsDeleteUser(userSlug string) error {
	if req, err := http.NewRequest("DELETE", H.Conf.GO_FIBER_VAULTS_URL + "/users/" + userSlug, nil);
	err != nil {
		return err
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		return err
	} else if resp.StatusCode != fiber.StatusNoContent {
		if respBody, err := io.ReadAll(resp.Body); err != nil {
			return errors.New("vaultsDeleteUser error, couldn't read response")
		} else {
			return errors.New(string(respBody))
		}
	}

	return nil
}
