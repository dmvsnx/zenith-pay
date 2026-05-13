package seed

import (
	"os"

	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
	"github.com/savanyv/zenith-pay/internal/utils/logger"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB, bc helpers.BcryptHelper) {
	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	email := os.Getenv("ADMIN_EMAIL")
	fullName := os.Getenv("ADMIN_FULL_NAME")

	if username == "" || password == "" {
		logger.Log.Info().Msg("ADMIN seed skipped (env not set)")
		return
	}

	hashedPassword, err := bc.HashPassword(password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to hash password")
		return
	}

	var admin model.User
	err = db.Where("username = ?", username).First(&admin).Error
	if err == nil {
		logger.Log.Info().Msg("ADMIN already exists, skipping seed")
		return
	}

	if err != gorm.ErrRecordNotFound {
		logger.Log.Error().Err(err).Msg("Failed to check admin existence")
		return
	}

	admin = model.User{
		Username: username,
		Password: hashedPassword,
		FullName: fullName,
		Email:    email,
		Role:     model.AdminRole,
		IsActive: true,
	}

	if err := db.Create(&admin).Error; err != nil {
		logger.Log.Error().Err(err).Msg("failed to seed admin")
		return
	}

	logger.Log.Info().Msg("ADMIN seeded successfully")
}
