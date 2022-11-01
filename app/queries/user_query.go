package queries

import (
	"errors"
	"github.com/aryanicosa/go-fiber-rest-api/app/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserQueries struct for queries from User model.
type UserQueries struct {
	DB *gorm.DB
}

// GetUserByID query for getting one User by given ID.
func (q *UserQueries) GetUserByID(id uuid.UUID) (models.User, error) {
	// Define User variable.
	user := models.User{}

	// Send query to database.
	err := q.DB.Table("users", q.DB.Model(&user)).Where("id = ?", id).Find(&user).Error
	if err != nil {
		// Return empty object and error.
		return user, errors.New("unable get user, DB error")
	}

	// Return query result.
	return user, nil
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserQueries) GetUserByEmail(email string) (models.User, error) {
	// Define User variable.
	user := models.User{}

	// Send query to database.
	err := q.DB.Table("users", q.DB.Model(&user)).Where("email = ?", email).Find(&user).Error
	if err != nil {
		// Return empty object and error.
		return user, errors.New("unable get user, DB error")
	}

	// Return query result.
	return user, nil
}

// CreateUser query for creating a new user by given email and password hash.
func (q *UserQueries) CreateUser(u *models.User) error {
	// Send query to database.
	err := q.DB.Table("users").Create(&models.User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		UserRole:     u.UserRole,
		UserStatus:   u.UserStatus,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
	}).Error
	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
