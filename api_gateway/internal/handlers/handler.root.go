package handlers

import (
	"api_gateway/internal/datatransfers"

	"github.com/gofiber/fiber/v2"
)

func Root(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("API online!!!", nil))
}

func Healthy(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("API healthy!!!", nil))
}
