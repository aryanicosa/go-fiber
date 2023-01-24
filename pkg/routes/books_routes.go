package routes

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/controllers"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// BooksRoutes func for describe group of private routes.
func BooksRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/v1")

	// Routes for POST method:
	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook) // create a new book

	// Routes for PUT method:
	route.Put("/book/:id", middleware.JWTProtected(), controllers.UpdateBook) // update one book by ID

	// Routes for DELETE method:
	route.Delete("/book/:id", middleware.JWTProtected(), controllers.DeleteBook) // delete one book by ID

	// Routes for GET method:
	route.Get("/books", middleware.BasicAuth(), controllers.GetBooks)   // get list of all books
	route.Get("/book/:id", middleware.BasicAuth(), controllers.GetBook) // get one book by ID
}
