package utils

import "github.com/gofiber/fiber/v2"

func RespondWithError(
	c *fiber.Ctx, statusCode int, detail string, fieldErrors map[string][]string,
	nonFieldErrors []string,
) error {
	return c.Status(statusCode).JSON(&ErrorResponseBody{
		Detail:         detail,
		FieldErrors:    fieldErrors,
		NonFieldErrors: nonFieldErrors,
	})
}
