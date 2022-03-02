package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type MagicTokens struct {
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId  uuid.UUID `json:"user_id"`
	Token   string    `json:"token"`
	Website string    `json:"website"`
}

func MigrateMagicTokens(db *gorm.DB) error {
	err := db.AutoMigrate(&MagicTokens{})
	return err
}
