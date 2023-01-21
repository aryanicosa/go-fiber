package middleware

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/response"
	"github.com/gofiber/fiber/v2"
	jwtMiddleware "github.com/gofiber/jwt/v2"
	"os"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 400 and failed bad request error.
	if err.Error() == "Missing or malformed JWT" {
		return response.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	// Return status 401 and failed authentication error.
	return response.RespondError(c, fiber.StatusUnauthorized, err.Error())
}
