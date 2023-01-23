package middleware

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	jwtMiddleware "github.com/gofiber/jwt/v2"
	"os"
)

func BasicAuth() func(*fiber.Ctx) error {
	config := basicauth.Config{
		Users: map[string]string{
			os.Getenv("BASIC_AUTH_USER"): os.Getenv("BASIC_AUTH_PASSWORD"),
		},
		Realm: "Forbidden",
		Authorizer: func(user, pass string) bool {
			if user == os.Getenv("BASIC_AUTH_USER") && pass == os.Getenv("BASIC_AUTH_PASSWORD") {
				return true
			}
			return false
		},
		Unauthorized: func(c *fiber.Ctx) error {
			return response.RespondError(c, fiber.StatusUnauthorized, "unauthorize access")
		},
		ContextUsername: "_user",
		ContextPassword: "_pass",
	}
	return basicauth.New(config)
}

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
