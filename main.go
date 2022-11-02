package main

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/configs"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/middleware"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/routes"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"github.com/aryanicosa/go-fiber-rest-api/platform/migrations"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// define a new fiber app with config
	app := fiber.New(config)

	// middlewares
	middleware.FiberMiddleware(app)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// migration
	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
	err = migrations.Migrate(migrationFileSource)
	if err != nil {
		log.Fatal("database migration fail")
	}

	// connect to db
	_, err = database.SqlConnection()
	if err != nil {
		log.Fatal("could not load database")
	}

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.UsersRoutes(app)   // Register a public routes for app.
	routes.BooksRoutes(app)   // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
