package database

import (
	"fmt"
	"github.com/aryanicosa/go-fiber-rest-api/app/queries"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Queries struct {
	*queries.UserQueries // load queries from User model
	*queries.BookQueries // load queries from Book model
}

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

func UserConn() (*queries.UserQueries, error) {
	// Configure any package-level settings
	return &queries.UserQueries{DB: db}, nil
}

func BookConn() (*queries.BookQueries, error) {
	// Configure any package-level settings
	return &queries.BookQueries{DB: db}, nil
}