package handlers

import (
	"api_gateway/internal/datatransfers"

	"github.com/gofiber/fiber/v2"
)

type CommonHandler struct{}

func NewCommonHandler() CommonHandler {
	return CommonHandler{}
}

func (h *CommonHandler) Root(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("API online!!!", nil))
}

func (h *CommonHandler) Healthy(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(datatransfers.ResponseSuccess("API healthy!!!", nil))
}
