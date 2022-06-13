package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Users struct {
	ID			uuid.UUID  	`json:"id" gorm:"type:uuid;primary_key;"`
	Name		string		`json:"name"`
	Email		string		`Json:"email" sql:"unique"`
	Address		string		`json:"address"`
	Password	[]byte		`json:"-"`
	CreatedAt 	*time.Time 	`json:"created_at,omitempty"`
	UpdatedAt 	*time.Time 	`json:"updated_at,omitempty"`
}

func MigrateUsers(db *gorm.DB) error {
    err := db.AutoMigrate(&Users{})
    return err
}