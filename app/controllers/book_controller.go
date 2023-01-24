package controllers

import (
	"github.com/aryanicosa/go-fiber-rest-api/pkg/response"
	"github.com/aryanicosa/go-fiber-rest-api/platform/database"
	"github.com/google/uuid"
	"time"

	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/repository"
	"github.com/aryanicosa/go-fiber-rest-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func GetBooks(c *fiber.Ctx) error {
	// Get all books.
	db := database.BookDB()
	books, err := db.GetBooks()
	if err != nil {
		// Return, if books not found.
		return response.RespondError(c, fiber.StatusNotFound, "books were not found")
	}

	var booksForPublic []models.BookForPublic
	for _, book := range books {
		bookForPublic := models.BookForPublic{
			ID:         book.ID,
			Title:      book.Title,
			Author:     book.Author,
			BookStatus: book.BookStatus,
			BookAttrs:  book.BookAttrs,
		}
		booksForPublic = append(booksForPublic, bookForPublic)
	}

	allBooks := &models.AllBooks{
		Books: booksForPublic,
		Count: int64(len(books)),
	}
	// Return status 200 OK.
	return response.RespondSuccess(c, fiber.StatusOK, allBooks)
}

func GetBook(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get book by ID.
	db := database.BookDB()
	book, err := db.GetBookById(id)
	if err != nil {
		// Return, if book not found.
		return response.RespondError(c, fiber.StatusNotFound, "book with the given ID is not found")
	}

	// Return status 200 OK.
	return response.RespondSuccess(c, fiber.StatusOK, book)
}

func CreateBook(c *fiber.Ctx) error {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusNotFound, err.Error())
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return response.RespondError(c, fiber.StatusUnauthorized, "unauthorized or expired token")
	}

	// Set credential `book:create` from JWT data of current book.
	credential := claims.Credentials[repository.BookCreateCredential]

	// Only user with `book:create` credential can create a new book.
	if !credential {
		// Return status 403 and permission denied error message.
		return response.RespondError(c, fiber.StatusForbidden, "permission denied, credential not eligible")
	}

	// Create new Book struct
	book := &models.Book{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, "unable to parse request body")
	}

	// Create a new validator for a Book model.
	validate := utils.NewValidator()

	// Set initialized default data for book:
	book.ID = uuid.New()
	book.CreatedAt = time.Now()
	book.UserID = claims.UserID
	book.BookStatus = 1 // 0 == draft, 1 == active

	// Validate book fields.
	if err := validate.Struct(book); err != nil {
		// Return, if some fields are not valid.
		return response.RespondError(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
	}

	// Create book by given model.
	db := database.BookDB()
	if err := db.CreateBook(book); err != nil {
		// Return status 500 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	return response.RespondSuccess(c, fiber.StatusCreated, book)
}

func UpdateBook(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return response.RespondError(c, fiber.StatusUnauthorized, "unauthorized or expired token")
	}

	// Set credential `book:update` from JWT data of current book.
	credential := claims.Credentials[repository.BookUpdateCredential]

	// Only book creator with `book:update` credential can update his book.
	if !credential {
		// Return status 403 and permission denied error message.
		return response.RespondError(c, fiber.StatusForbidden, "permission denied, credential not eligible")
	}

	// Create new Book struct
	book := &models.Book{}

	// Check, if received JSON data is valid.
	if err := c.BodyParser(book); err != nil {
		// Return status 400 and error message.
		return response.RespondError(c, fiber.StatusBadRequest, err.Error())
	}

	// Checking, if book with given ID is exists.
	db := database.BookDB()
	foundedBook, err := db.GetBookById(id)
	if err != nil {
		// Return status 404 and book not found error.
		return response.RespondError(c, fiber.StatusNotFound, "book with given ID not found")
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his book.
	if foundedBook.UserID == userID {
		// Set initialized default data for book:
		book.UpdatedAt = time.Now()

		// Create a new validator for a Book model.
		validate := utils.NewValidator()

		// Validate book fields.
		if err := validate.Struct(book); err != nil {
			// Return, if some fields are not valid.
			return response.RespondError(c, fiber.StatusBadRequest, utils.ValidatorErrors(err))
		}

		// Update book by given ID.
		if err := db.UpdateBook(foundedBook.ID, book); err != nil {
			// Return status 500 and error message.
			return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
		}

		// Return status 201.
		return response.RespondSuccess(c, fiber.StatusCreated, book)
	} else {
		// Return status 403 and permission denied error message.
		return response.RespondError(c, fiber.StatusForbidden, "permission denied, only the creator can update this book")
	}
}

func DeleteBook(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return response.RespondError(c, fiber.StatusUnauthorized, "unauthorized or expired token")
	}

	// Set credential `book:delete` from JWT data of current book.
	credential := claims.Credentials[repository.BookDeleteCredential]

	// Only book creator with `book:delete` credential can delete his book.
	if !credential {
		// Return status 403 and permission denied error message.
		return response.RespondError(c, fiber.StatusForbidden, "permission denied, credential not eligible")
	}

	// Checking, if book with given ID is exists.
	db := database.BookDB()
	foundedBook, err := db.GetBookById(id)
	if err != nil {
		// Return status 404 and book not found error.
		return response.RespondError(c, fiber.StatusNotFound, "book with given ID not found")
	}

	// Set user ID from JWT data of current user.
	userID := claims.UserID

	// Only the creator can delete his book.
	if foundedBook.UserID == userID {
		// Delete book by given ID.
		if err := db.DeleteBook(foundedBook.ID); err != nil {
			// Return status 500 and error message.
			return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
		}

		// Return status 204 no content.
		return response.RespondSuccess(c, fiber.StatusNoContent, "")
	} else {
		// Return status 403 and permission denied error message.
		return response.RespondError(c, fiber.StatusForbidden, "permission denied, only the creator can delete this book")
	}
}
