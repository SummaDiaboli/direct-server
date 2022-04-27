package models

import "gorm.io/gorm"

// Processes the migration functions for the different models
func MigrateTables(db *gorm.DB) error {
	err := MigrateUsers(db)
	if err != nil {
		return err
	}

	err = MigrateWebsites(db)
	if err != nil {
		return err
	}

	err = MigrateAuthTokens(db)
	if err != nil {
		return err
	}

	err = MigrateAuthedWebsites(db)
	if err != nil {
		return err
	}

	// err = MigrateQRCodes(db)
	// if err != nil {
	// 	return err
	// }

	return err
}
