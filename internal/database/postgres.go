package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/savanyv/zenith-pay/config"
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel: logger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful: true,
		},
	)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	if cfg.AppEnv == "development" {
		if err := db.AutoMigrate(
			&model.User{},
			&model.Category{},
			&model.Product{},
			&model.Transaction{},
			&model.TransactionItems{},
			&model.Shift{},
		); err != nil {
			log.Fatal(err)
		}
	}

	DB = db
	log.Println("Database connected")
	return db, nil
}
