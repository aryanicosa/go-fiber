package queries

import (
	"errors"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BookQueries struct for queries from Book model.
type BookQueries struct {
	*gorm.DB
}

// CreateBook method for creating book by given Book object.
func (q *BookQueries) CreateBook(b *models.Book) error {
	// Send query to database.
	err := q.DB.Table("books").Create(&models.Book{
		ID:         b.ID,
		CreatedAt:  b.CreatedAt,
		UpdatedAt:  b.UpdatedAt,
		UserID:     b.UserID,
		Title:      b.Title,
		Author:     b.Author,
		BookStatus: b.BookStatus,
		BookAttrs:  b.BookAttrs,
	}).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// GetBooks method for getting all books.
func (q *BookQueries) GetBooks() ([]*models.Book, error) {
	// Define books variable.
	books := []*models.Book{}

	// Send query to database.
	err := q.DB.Table("books").Find(&books).Error
	if err != nil {
		// Return empty object and error.
		return nil, err
	}

	// Return query result.
	return books, nil
}

// GetBookById method for getting one book by given ID.
func (q *BookQueries) GetBookById(id uuid.UUID) (models.Book, error) {
	// Define book variable.
	book := models.Book{}

	// Send query to database.
	err := q.DB.Table("books", q.DB.Model(&book)).Where("id = ?", id).Find(&book).Error
	if err != nil {
		// Return empty object and error.
		return book, errors.New("unable get user, DB error")
	}

	// Return query result.
	return book, nil
}

// GetBooksByAuthor method for getting all books by given author.
func (q *BookQueries) GetBooksByAuthor(author string) ([]models.Book, error) {
	// Define books variable.
	books := []models.Book{}

	// Send query to database.
	err := q.DB.Table("books", q.DB.Model(&books)).Where("author = ?", author).Find(&books).Error
	if err != nil {
		// Return empty object and error.
		return books, errors.New("unable get user, DB error")
	}

	// Return query result.
	return books, nil
}

// UpdateBook method for updating book by given Book object.
func (q *BookQueries) UpdateBook(id uuid.UUID, b *models.Book) error {
	// Send query to database.
	err := q.DB.Table("books").Where("id = ?", id).Updates(b).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteBook method for delete book by given ID.
func (q *BookQueries) DeleteBook(id uuid.UUID) error {
	// Send query to database.
	err := q.Table("books").Where("id = ?", id).Delete(&models.Book{}).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
