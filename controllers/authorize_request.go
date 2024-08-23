package controllers

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
	"github.com/liobrdev/simplepasswords_api_gateway/utils"
)

type userContextKey struct{}

func (H Handler) AuthorizeRequest(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	// Check is null or ""
	if utils.TokenNullRegexp.Match([]byte(authHeader)) {
		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	// Check Authorization header format
	if !utils.AuthHeaderRegexp.Match([]byte(authHeader)) {
		H.logger(c, c.Get("Client-Operation"), "Authorization", authHeader, "error", utils.ErrorToken)

		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	authToken := authHeader[6:]

	var thisSession models.ClientSession

	// Query session from token provided in header
	if result := H.DBs.ApiGateway.Preload("User").Where("token_key = ?", authToken[:16]).Limit(1).
	Find(&thisSession); result.Error != nil {
		H.logger(c, c.Get("Client-Operation"), result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if n := result.RowsAffected; n == 0 {
		// No matches - client can try again or re-login
		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	} else if n != 1 {
		H.logger(
			c, c.Get("Client-Operation"), "result.RowsAffected", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	// Check session user is active
	if !thisSession.User.IsActive {
		H.logger(
			c, c.Get("Client-Operation"), "!thisSession.User.IsActive",
			"token_key = " + authToken[:16] + " ; user_slug = " + thisSession.UserSlug, "error",
			utils.ErrorBadClient,
		)

		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	// Check request IP address matches requested session
	if c.IP() != thisSession.ClientIP {
		H.logger(
			c, c.Get("Client-Operation"), "c.IP() != thisSession.ClientIP",
			c.IP() + " != " + thisSession.ClientIP, "error", utils.ErrorIPMismatch,
		)

		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	var userAllSessions []models.ClientSession

	// Get all user's sessions for cleanup
	if result := H.DBs.ApiGateway.Where("user_slug = ?", thisSession.UserSlug).Find(&userAllSessions);
	result.Error != nil {
		H.logger(c, c.Get("Client-Operation"), result.Error.Error(), "", "error", utils.ErrorFailedDB)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	} else if n := result.RowsAffected; n < 1 {
		H.logger(
			c, c.Get("Client-Operation"), "result.RowsAffected < 1", strconv.FormatInt(n, 10), "error",
			utils.ErrorFailedDB,
		)

		return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
	}

	thisSessionExpired := false
	now := time.Now().UTC()
	
	// Clean up expired sessions
	for _, session := range userAllSessions {
		if session.ExpiresAt.Before(now) {
			if session.TokenKey == thisSession.TokenKey {
				thisSessionExpired = true
			}

			if result := H.DBs.ApiGateway.Delete(&session); result.Error != nil {
				H.logger(
					c, c.Get("Client-Operation"), result.Error.Error(), "", "error", utils.ErrorFailedDB,
				)

				return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
			} else if n := result.RowsAffected; n != 1 {
				H.logger(
					c, c.Get("Client-Operation"), "result.RowsAffected != 1", strconv.FormatInt(n, 10),
					"error", utils.ErrorFailedDB,
				)

				return utils.RespondWithError(c, 500, utils.ErrorServer, nil, nil)
			}
		}
	}

	// Client can try again if session expired
	if thisSessionExpired {
		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	if bytes.Equal(utils.HashToken(authToken), thisSession.Digest) {
		newExpiresAt := now.Add(time.Duration(15) * time.Minute)

		// Throttle updates to 'expires_at' by 60 sec
		if newExpiresAt.Sub(thisSession.ExpiresAt).Seconds() > 60 {
			thisSession.ExpiresAt = newExpiresAt
			H.DBs.ApiGateway.Save(&thisSession)
		}
	} else {
		H.logger(
			c, c.Get("Client-Operation"), "utils.HashToken(authToken) != thisSession.Digest",
			"authToken = " + authToken, "error", utils.ErrorToken,
		)
	
		return utils.RespondWithError(c, 401, utils.ErrorToken, nil, nil)
	}

	c.SetUserContext(context.WithValue(context.Background(), userContextKey{}, &thisSession.User))

	return c.Next()
}
