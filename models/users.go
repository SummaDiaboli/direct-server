package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// The users that create accounts for the service
type Users struct {
	ID            uuid.UUID      `gorm:"primary_key;type:uuid;default:uuid_generate_v4();" json:"id"`
	Username      string         `gorm:"unique" json:"username"`
	Email         string         `json:"email"`
	Created       datatypes.Date `json:"created"`
	TokenDuration int            `json:"token_duration" gorm:"default:7"`
	// Website  string    `json:"website"`
}

// Create users table
func MigrateUsers(db *gorm.DB) error {
	err := db.AutoMigrate(&Users{})
	return err
}
