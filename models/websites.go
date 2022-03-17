package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// The website the users current have active so that they can maintain control on their login status
type Websites struct {
	ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	// Url           string    `json:"url"`
	// Name  string `json:"website"`
	Token string `json:"token"`
	// Expires       string    `json:"expires"`
	UserId        uuid.UUID `json:"user_id"`
	Authenticated bool      `json:"authenticated"`
}

// Create websites table
func MigrateWebsites(db *gorm.DB) error {
	err := db.AutoMigrate(&Websites{})
	return err
}
