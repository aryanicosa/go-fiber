package database

import (
	"fmt"
	"github.com/aryanicosa/go-fiber-rest-api/app/queries"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
	TimeZone string
}

type Queries struct {
	*queries.UserQueries // load queries from User model
	*queries.BookQueries // load queries from Book model
}

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*Queries, error) {
	// Build PostgreSQL connection URL.
	postgresConnURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgreSQL.
	db, err := gorm.Open(postgres.Open(postgresConnURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	return &Queries{
		// Set queries from models:
		UserQueries: &queries.UserQueries{DB: db}, // from User model
		BookQueries: &queries.BookQueries{DB: db}, // from Book model
	}, nil
}
