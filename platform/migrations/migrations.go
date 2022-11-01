package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(source string) error {
	sourceStr := fmt.Sprintf("file://%s", source)
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
	db, err := sql.Open("postgres", dbUrl)
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(sourceStr, "postgres", driver)
	if err != nil {
		log.Fatal(err)
		return err
	}
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	return nil
}
