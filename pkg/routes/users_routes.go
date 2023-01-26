package routes

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/controllers"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// UsersRoutes func for describe group of public routes.
func UsersRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for public routes:
	route.Post("/user/sign/up", middleware.BasicAuth(), controllers.UserSignUp) // register a new user
	route.Post("/user/sign/in", middleware.BasicAuth(), controllers.UserSignIn) // auth, return AccessToken & RefreshToken tokens

	// Routes for privates routes:
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)   // de-authorization user
	route.Post("/user/sign/renew", middleware.JWTProtected(), controllers.RenewTokens) // renew AccessToken & RefreshToken tokens

}
