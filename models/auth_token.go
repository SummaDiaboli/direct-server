package models

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type AuthTokens struct {
	ID     uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserId uuid.UUID ` gorm:"unique" json:"user_id"`
	Token  string    `json:"token"`
}

func MigrateAuthTokens(db *gorm.DB) error {
	err := db.AutoMigrate(&AuthTokens{})
	return err
}
