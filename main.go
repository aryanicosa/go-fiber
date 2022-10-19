package main

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/configs"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/middleware"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/routes"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
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

	// Routes.
	routes.SwaggerRoute(app)  // Register a route for API Docs (Swagger).
	routes.PublicRoutes(app)  // Register a public routes for app.
	routes.PrivateRoutes(app) // Register a private routes for app.
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
