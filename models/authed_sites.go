package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// The website the users current have active so that they can maintain control on their login status
type AuthedWebsites struct {
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	Token   string    `json:"token"`
	Referer string    `json:"referer"`
	UserId  uuid.UUID `json:"user_id"`
	// Expired bool           `json:"expired"`
	Created datatypes.Date `json:"created"`
	Expires datatypes.Date `json:"expires"`
}

// Create websites table
func MigrateAuthedWebsites(db *gorm.DB) error {
	err := db.AutoMigrate(&AuthedWebsites{})
	return err
}
