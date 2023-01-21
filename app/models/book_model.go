package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

// AllBooks struct to return all books
type AllBooks struct {
	Books []Book `json:"books"`
	Count int64
}

// Book struct to describe book object.
type Book struct {
	ID         uuid.UUID `json:"id" validate:"required,uuid"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     uuid.UUID `json:"user_id" validate:"required,uuid"`
	Title      string    `json:"title" validate:"required,lte=255"`
	Author     string    `json:"author" validate:"required,lte=255"`
	BookStatus int       `json:"book_status" validate:"required,len=1"`
	BookAttrs  BookAttrs `json:"book_attrs" validate:"required,dive"`
}

// BookAttrs struct to describe book attributes.
type BookAttrs struct {
	Picture     string `json:"picture"`
	Description string `json:"description"`
	Rating      int    `json:"rating" validate:"min=0,max=10"`
}

// Value make the BookAttrs struct implement the driver.Valuer interface.
// This method simply returns the JSON-encoded representation of the struct.
func (b BookAttrs) Value() (driver.Value, error) {
	return json.Marshal(b)
}

// Scan make the BookAttrs struct implement the sql.Scanner interface.
// This method simply decodes a JSON-encoded value into the struct fields.
func (b *BookAttrs) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(j, &b)
}
