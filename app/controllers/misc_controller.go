package controllers

import (
	"encoding/base64"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// EncodeBase64 godoc
// @Summary 	Encode String to Base64
// @Description Encode input string to Base64 string
// @Accept 		json
// @Produce 	json
// @Tags 						Miscellaneous
// @Param StringToEncode body string true "arbitrary string"
// @Success 200 {string} string
// @Failure 400 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/misc/base64encode [post]
func EncodeBase64(c *fiber.Ctx) error {
	EncodedString := base64.StdEncoding.EncodeToString(c.Body())
	// Return status 200 OK.
	return response.RespondSuccess(c, fiber.StatusCreated, EncodedString)
}
