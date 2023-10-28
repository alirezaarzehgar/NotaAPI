package database

import (
	"gorm.io/gorm"

	"github.com/Asrez/NotaAPI/config"
	"github.com/Asrez/NotaAPI/models"
	"github.com/Asrez/NotaAPI/utils"
)

func usersMigrate(db *gorm.DB) error {
	c := config.Admin()
	admin := models.User{
		Username: c.Username,
		Email:    c.Email,
		Password: utils.HashPassword(c.Password),
		Role:     models.USERS_ROLE_ADMIN,
	}
	return db.Create(&admin).Error
}

func Migrate(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		return nil
	}

	err := db.AutoMigrate(&models.User{}, &models.Token{}, &models.Guest{}, &models.Story{})
	if err != nil {
		return err
	}

	if err := usersMigrate(db); err != nil {
		return err
	}

	return nil
}
