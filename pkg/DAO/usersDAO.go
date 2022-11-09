package DAO

import (
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/google/uuid"
)

type UserQueries interface {
	GetUserByID(id uuid.UUID) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	CreateUser(u *models.User) error
}
