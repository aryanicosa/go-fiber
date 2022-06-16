package user

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	jwtware "github.com/gofiber/jwt/v3"
)

func JwtAuth() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("SECRET_KEY")),
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func BasicAuth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("BASIC_AUTH_USERNAME"): os.Getenv("BASIC_AUTH_PASSWORD"),
		},
	})
}

func (s *Service) SetupRoutes(app *fiber.App) {
	api := app.Group("/user")
	api.Post("/", BasicAuth(), s.CreateUser)
	api.Post("/login", BasicAuth(), s.LoginUser)
	api.Post("/logout", JwtAuth(), s.LogoutUser)
	api.Get("/all", JwtAuth(), s.GetUsers)
	api.Get("/:id", JwtAuth(), s.GetUser)
	api.Put("/:id", JwtAuth(), s.UpdateUser)
	api.Delete("/:id", JwtAuth(), s.DeleteUser)

}
