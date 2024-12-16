package handler

import "github.com/gofiber/fiber/v2"

func Healthz(c *fiber.Ctx) error {
	return c.SendStatus(200)
}
