package routes

import (
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"github.com/aryanicosa/go-fiber-rest-api/platform/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"
)

var AppTest *fiber.App

func TestMain(m *testing.M) {
	// Load .env.test file from the root folder.
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatal(err)
	}

	// Define Fiber AppTest.
	AppTest = fiber.New()

	// init connect to db
	_, err := database.InitDBConnection()
	if err != nil {
		log.Fatal("fail to load database")
	}

	// migration
	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
	err = migrations.Migrate(migrationFileSource)
	if err != nil {
		log.Fatal("database migration fail")
	}

	// Define routes.
	BooksRoutes(AppTest)
	UsersRoutes(AppTest)
	MiscRoutes(AppTest)

	os.Exit(m.Run())
}
