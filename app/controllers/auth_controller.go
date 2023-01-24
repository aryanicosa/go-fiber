package controllers

import (
	"context"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/response"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"github.com/aryanicosa/go-fiber-rest-api/platform/cache"
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UserSignUp(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signUp := &models.SignUp{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signUp); err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, "unable to parse request body")
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(signUp); err != nil {
		// Return, if some fields are not valid.
		return response.RespondError(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Checking role from sign up data.
	role, err := utils.VerifyRole(signUp.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	// Create a new user struct.
	user := &models.User{}

	// Set initialized default data for user:
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.Email = signUp.Email
	user.PasswordHash = utils.GeneratePassword(signUp.Password)
	user.UserStatus = 1 // 0 == blocked, 1 == active
	user.UserRole = role

	// Validate user fields.
	if err := validate.Struct(user); err != nil {
		// Return, if some fields are not valid.
		return response.RespondError(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Create a new user with validated data.
	db := database.UserDB()
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}
	if err := db.CreateUser(user); err != nil {
		// Return status 500 and create user process error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Delete password hash field from JSON view.
	user.PasswordHash = ""

	// Return status 201 OK.
	return response.RespondSuccess(c, fiber.StatusCreated, user)
}

// UserSignIn godoc
// @Summary 	Sign In
// @Description Sign In a User to get access token
// @Description value of Authorization field is "Basic base64string_of_username:secret"
// @Description use /v1/misc/encode to generate your base64 encoded string
// @Accept 		json
// @Produce 	json
// @Tags 						User
// @Param Authorization header string true "Basic Auth"
// @Param models.SignIn body models.SignIn true "User Credentials"
// @Success 200 {object} utils.Tokens
// @Failure 400 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/user/sign/in [post]
func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &models.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, "unable to parse request body")
	}

	// Get user by email.
	db := database.UserDB()
	foundedUser, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		// Return, if user not found.
		return response.RespondError(c, fiber.StatusNotFound, err.Error())
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(foundedUser.PasswordHash, signIn.Password)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return response.RespondError(c, fiber.StatusUnauthorized, "wrong user email address or password")
	}

	// Get role credentials from founded user.
	credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
	if err != nil {
		// Return status 403 and error message.
		return response.RespondError(c, fiber.StatusForbidden, err.Error())
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(foundedUser.ID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Define user ID.
	userID := foundedUser.ID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), userID, tokens.RefreshToken, 0).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return response.RespondError(c, fiber.StatusInternalServerError, errSaveToRedis.Error())
	}

	// Return status 200 OK.
	return response.RespondSuccess(c, fiber.StatusOK, tokens)
}

func UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return response.RespondError(c, fiber.StatusInternalServerError, errDelFromRedis.Error())
	}

	// Return status 204 no content.
	return response.RespondSuccess(c, fiber.StatusNoContent, "")
}

func RenewTokens(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// Checking, if now time greater than AccessToken token expiration time.
	if now > expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return response.RespondError(c, fiber.StatusUnauthorized, "unauthorized or expired token")
	}

	// Create a new renewal refresh token struct.
	renew := &models.Renew{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renew); err != nil {
		// Return, if JSON data is not correct.
		return response.RespondError(c, fiber.StatusBadRequest, "unable to parse request body")
	}

	// Set expiration time from RefreshToken token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	// Checking, if now time greater than RefreshToken token expiration time.
	if now < expiresRefreshToken {
		// Define user ID.
		userID := claims.UserID

		// Get user by ID.
		db := database.UserDB()
		foundedUser, err := db.GetUserByID(userID)
		if err != nil {
			// Return, if user not found.
			return response.RespondError(c, fiber.StatusNotFound, "user with the given ID is not found")
		}

		// Get role credentials from founded user.
		credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
		if err != nil {
			// Return status 400 and error message.
			return response.RespondError(c, fiber.StatusBadRequest, err.Error())
		}

		// Generate JWT AccessToken & RefreshToken tokens.
		tokens, err := utils.GenerateNewTokens(userID.String(), credentials)
		if err != nil {
			// Return status 500 and token generation error.
			return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
		}

		// Create a new Redis connection.
		connRedis, err := cache.RedisConnection()
		if err != nil {
			// Return status 500 and Redis connection error.
			return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
		}

		// Save refresh token to Redis.
		errRedis := connRedis.Set(context.Background(), userID.String(), tokens.RefreshToken, 0).Err()
		if errRedis != nil {
			// Return status 500 and Redis connection error.
			return response.RespondError(c, fiber.StatusInternalServerError, errRedis.Error())
		}

		return response.RespondSuccess(c, fiber.StatusOK, tokens)
	} else {
		// Return status 401 and unauthorized error message.
		return response.RespondError(c, fiber.StatusUnauthorized, "unauthorized, your session was ended earlier")
	}
}
