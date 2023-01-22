package response

import "github.com/gofiber/fiber/v2"

// HTTPError struct to HTTPError object.
type HTTPError struct {
	ErrorMessage interface{} `json:"errorMessage"`
}

func RespondError(c *fiber.Ctx, responseCode int, errMessage interface{}) error {
	errorJson := &HTTPError{
		ErrorMessage: errMessage,
	}
	return c.Status(responseCode).JSON(errorJson)
}

func RespondSuccess(c *fiber.Ctx, responseCode int, data interface{}) error {
	return c.Status(responseCode).JSON(data)
}
