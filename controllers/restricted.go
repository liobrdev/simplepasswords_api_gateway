package controllers

import "github.com/gofiber/fiber/v2"

func (H Handler) Restricted(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
