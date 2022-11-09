package DAO

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/google/uuid"
)

type BookQueries interface {
	CreateBook(b *models.Book) error
	GetBooks() ([]models.Book, error)
	GetBookById(id uuid.UUID) (models.Book, error)
	GetBooksByAuthor(author string) ([]models.Book, error)
	UpdateBook(id uuid.UUID, b *models.Book) error
	DeleteBook(id uuid.UUID) error
}
