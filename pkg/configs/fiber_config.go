package configs

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"time"
)

func FiberConfig() fiber.Config {
	// define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))

	// return fiber configuration
	return fiber.Config{
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
