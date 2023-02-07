package migrations

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func Migrate(source string) error {
	sourceStr := fmt.Sprintf("file://%s", source) // "///" means absolute path "//" mean relative path
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_NAME"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("fail open db connection")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("fail to get db driver")
	}
	m, err := migrate.NewWithDatabaseInstance(sourceStr, "postgres", driver)
	if err != nil {
		log.Fatal(err)
		return err
	}
	//nolint:golint // make it catch error once no changes
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No changes, database schema is up to date")
			return nil
		}
		log.Fatal(err)
		return err
	} // or m.Step(2) if you want to explicitly set the number of migrations to run
	log.Println("database schema migration done")
	return nil
}
