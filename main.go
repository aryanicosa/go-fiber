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

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Go Fiber Rest API
// @contact.name API Support
// @contact.email aryanicosa@gmail.com
// @BasePath /
// @securityDefinitions.basic BasicAuth
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	// init connect to db
	_, errInitDb := database.InitDBConnection()
	if errInitDb != nil {
		log.Fatal("could not load database")
	}

	// migration
	migrationFileSource := os.Getenv("SQL_SOURCE_PATH")
	err = migrations.Migrate(migrationFileSource)
	if err != nil {
		log.Fatal("database migration fail")
	}

	// Routes.
	routes.SwaggerRoute(app) // Register a route for API Docs (Swagger).
	routes.UsersRoutes(app)
	routes.BooksRoutes(app)
	routes.MiscRoutes(app)
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
