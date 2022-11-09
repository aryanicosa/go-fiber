package controllers

import (
	"context"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Create a new validator for a User model.
	validate := utils.NewValidator()

	// Validate sign up fields.
	if err := validate.Struct(signUp); err != nil {
		// Return, if some fields are not valid.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Checking role from sign up data.
	role, err := utils.VerifyRole(signUp.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   utils.ValidatorErrors(err),
		})
	}

	// Create a new user with validated data.
	db, err := database.UserDB()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}
	if err := db.CreateUser(user); err != nil {
		// Return status 500 and create user process error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Delete password hash field from JSON view.
	user.PasswordHash = ""

	// Return status 201 OK.
	return c.Status(fiber.StatusCreated).JSON(user)
}

func UserSignIn(c *fiber.Ctx) error {
	// Create a new user auth struct.
	signIn := &models.SignIn{}

	// Checking received data from JSON body.
	if err := c.BodyParser(signIn); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Get user by email.
	db, err := database.UserDB()
	foundedUser, err := db.GetUserByEmail(signIn.Email)
	if err != nil {
		// Return, if user not found.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Compare given user password with stored in found user.
	compareUserPassword := utils.ComparePasswords(foundedUser.PasswordHash, signIn.Password)
	if !compareUserPassword {
		// Return, if password is not compare to stored in database.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   "wrong user email address or password",
		})
	}

	// Get role credentials from founded user.
	credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(foundedUser.ID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userID := foundedUser.ID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), userID, tokens.Refresh, 0).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errSaveToRedis.Error(),
		})
	}

	// Return status 200 OK.
	return c.Status(fiber.StatusOK).JSON(
		utils.Tokens{
			Access:  tokens.Access,
			Refresh: tokens.Refresh,
		})
}

func UserSignOut(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create a new Redis connection.
	connRedis, err := cache.RedisConnection()
	if err != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   errDelFromRedis.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

func RenewTokens(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// Checking, if now time greather than Access token expiration time.
	if now > expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, check expiration time of your token",
		})
	}

	// Create a new renew refresh token struct.
	renew := &models.Renew{}

	// Checking received data from JSON body.
	if err := c.BodyParser(renew); err != nil {
		// Return, if JSON data is not correct.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now < expiresRefreshToken {
		// Define user ID.
		userID := claims.UserID

		// Get user by ID.
		db, err := database.UserDB()
		foundedUser, err := db.GetUserByID(userID)
		if err != nil {
			// Return, if user not found.
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": true,
				"msg":   "user with the given ID is not found",
			})
		}

		// Get role credentials from founded user.
		credentials, err := utils.GetCredentialsByRole(foundedUser.UserRole)
		if err != nil {
			// Return status 400 and error message.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Generate JWT Access & Refresh tokens.
		tokens, err := utils.GenerateNewTokens(userID.String(), credentials)
		if err != nil {
			// Return status 500 and token generation error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Create a new Redis connection.
		connRedis, err := cache.RedisConnection()
		if err != nil {
			// Return status 500 and Redis connection error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}

		// Save refresh token to Redis.
		errRedis := connRedis.Set(context.Background(), userID.String(), tokens.Refresh, 0).Err()
		if errRedis != nil {
			// Return status 500 and Redis connection error.
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   errRedis.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": false,
			"msg":   nil,
			"tokens": fiber.Map{
				"access":  tokens.Access,
				"refresh": tokens.Refresh,
			},
		})
	} else {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": true,
			"msg":   "unauthorized, your session was ended earlier",
		})
	}
}
