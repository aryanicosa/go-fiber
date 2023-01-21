package response

import "github.com/gofiber/fiber/v2"

func RespondError(c *fiber.Ctx, responseCode int, errMessage interface{}) error {
	return c.Status(responseCode).JSON(errMessage)
}

func RespondSuccess(c *fiber.Ctx, responseCode int, data interface{}) error {
	return c.Status(responseCode).JSON(data)
}
