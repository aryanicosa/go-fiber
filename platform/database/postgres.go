package database

import (
	"fmt"
	"github.com/aryanicosa/go-fiber-rest-api/app/queries"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

// InitDBConnection func for connection to PostgreSQL database.
func InitDBConnection() (*gorm.DB, error) {
	// Build PostgreSQL connection URL.
	postgresConnURL, err := utils.ConnectionURLBuilder("postgres")
	if err != nil {
		return nil, err
	}

	// Define database connection for PostgreSQL.
	db, err = gorm.Open(postgres.Open(postgresConnURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	return db, nil
}

// UserDB used for init users db query
func UserDB() (*queries.UserQueries, error) {
	return &queries.UserQueries{DB: db}, nil
}

// BookDB used for init books db query
func BookDB() (*queries.BookQueries, error) {
	return &queries.BookQueries{DB: db}, nil
}
