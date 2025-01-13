package database

import (
	"go-nginx/models"

	"gorm.io/gorm"
)

func Migrate(patentDB, userDB *gorm.DB) error {
	// Migrate patent database if provided
	if patentDB != nil {
		if err := patentDB.AutoMigrate(&models.Patent{}, &models.Result{}); err != nil {
			return err
		}
	}

	// Migrate user database if provided
	if userDB != nil {
		if err := userDB.AutoMigrate(&models.User{}); err != nil {
			return err
		}
	}

	return nil
}
