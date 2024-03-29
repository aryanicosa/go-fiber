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

// GetBooks godoc
// @Description Will display all books
// @Description Require Basic Auth
// @Summary Get All Books
// @Tags Book
// @Accept json
// @Produce json
// @Security BasicAuth
// @Success 200 {array} models.BookForPublic
// @Failure 404 {object} response.HTTPError
// @Router /v1/books [get]
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
		arrBookForPublic := []models.BookForPublic{
			{
				ID:         book.ID,
				Title:      book.Title,
				Author:     book.Author,
				BookStatus: book.BookStatus,
				BookAttrs:  book.BookAttrs,
			},
		}
		booksForPublic = arrBookForPublic
	}

	allBooks := &models.AllBooks{
		Books: booksForPublic,
		Count: int64(len(books)),
	}
	// Return status 200 OK.
	return response.RespondSuccess(c, fiber.StatusOK, allBooks)
}

// GetBook godoc
// @Description Will display specific book by it's ID
// @Description Require valid user token
// @Summary Get book by ID
// @Tags Book
// @Accept json
// @Produce json
// @Security BasicAuth
// @Param book_id path string true "Book ID"
// @Success 200 {object} models.Book
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/book/id [get]
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

// CreateBook godoc
// @Description Require valid user token
// @Summary Create new book
// @Tags Book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param models.Book body models.Book true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} response.HTTPError
// @Failure 403 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/book [post]
func CreateBook(c *fiber.Ctx) error {
	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusNotFound, err.Error())
	}

	// Set credential needed `book:delete` from JWT data of current book.
	credentialNeed := repository.BookCreateCredential

	if isError, errorCode, errorMessage := bookClaimCheck(claims, credentialNeed); isError {
		return response.RespondError(c, errorCode, errorMessage)
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

// UpdateBook godoc
// @Description Require valid user token
// @Summary Update a book
// @Tags Book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param book_id path string true "Book ID"
// @Param models.Book body models.Book true "Book data"
// @Success 201 {object} models.Book
// @Failure 400 {object} response.HTTPError
// @Failure 401 {object} response.HTTPError
// @Failure 403 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/book/id [put]
func UpdateBook(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Set credential needed `book:delete` from JWT data of current book.
	credentialNeed := repository.BookUpdateCredential

	if isError, errorCode, errorMessage := bookClaimCheck(claims, credentialNeed); isError {
		return response.RespondError(c, errorCode, errorMessage)
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

// DeleteBook godoc
// @Description Require valid user token
// @Summary Delete a book
// @Tags Book
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param book_id path string true "Book ID"
// @Success 204
// @Failure 400 {object} response.HTTPError
// @Failure 401 {object} response.HTTPError
// @Failure 403 {object} response.HTTPError
// @Failure 404 {object} response.HTTPError
// @Failure 500 {object} response.HTTPError
// @Router /v1/book/id [delete]
func DeleteBook(c *fiber.Ctx) error {
	// Catch book ID from URL.
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return response.RespondError(c, fiber.StatusInternalServerError, err.Error())
	}

	// Set credential needed `book:delete` from JWT data of current book.
	credentialNeed := repository.BookDeleteCredential

	if isError, errorCode, errorMessage := bookClaimCheck(claims, credentialNeed); isError {
		return response.RespondError(c, errorCode, errorMessage)
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

func bookClaimCheck(claims *utils.TokenMetadata, credentialNeeded string) (bool, int, interface{}) {
	now := time.Now().Unix()

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greater than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return true, fiber.StatusUnauthorized, "unauthorized or expired token"
	}

	isValidCredential := claims.Credentials[credentialNeeded]

	// Only book creator with `book:delete` credential can delete his book.
	if !isValidCredential {
		// Return status 403 and permission denied error message.
		return true, fiber.StatusForbidden, "permission denied, credential not eligible"
	}

	return false, 0, ""
}
