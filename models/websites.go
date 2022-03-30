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
	// Expires       string    `json:"expires"`
	Token         string    `json:"token"`
	Referer       string    `json:"referer"`
	UserId        uuid.UUID `gorm:"unique" json:"user_id"`
	Authenticated bool      `json:"authenticated"`
}

// Create websites table
func MigrateWebsites(db *gorm.DB) error {
	err := db.AutoMigrate(&Websites{})
	return err
}
