package controllers

import "github.com/gofiber/fiber/v2"

const RestrictedCaller string = "api_gateway/controllers/restricted.go"

func (H Handler) Restricted(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
