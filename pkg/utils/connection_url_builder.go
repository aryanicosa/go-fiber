package utils

import (
	"fmt"
	"os"
)

// ConnectionURLBuilder func for building url connection.
func ConnectionURLBuilder(str string) (string, error) {
	// define URL to connection
	var url string

	// switch given names.
	switch str {
	case "postgres":
		// url for postgre connection
		url = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_NAME"),
			os.Getenv("POSTGRES_SSL_MODE"),
			os.Getenv("POSTGRES_TIME_ZONE"),
		)
	case "redis":
		// url for redis connection
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		)
	case "fiber":
		// url for fiber connection
		url = fmt.Sprintf(
			"%s:%s",
			os.Getenv("SERVER_HOST"),
			os.Getenv("SERVER_PORT"),
		)
	default:
		// Return error message.
		return "", fmt.Errorf("connection name '%v' is not supported", str)
	}

	return url, nil
}
