package routes

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/controllers"
	"github.com/gofiber/fiber/v2"
)

// MiscRoutes func for handling non-functionality service.
func MiscRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for POST method:
	route.Post("/misc/base64encode", controllers.EncodeBase64) // encode input
}
