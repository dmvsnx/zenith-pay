package database

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/savanyv/zenith-pay/config"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/utils/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

type gormLogWriter struct {
	l zerolog.Logger
}

func (w *gormLogWriter) Printf(format string, args ...interface{}) {
	w.l.Warn().Msg(fmt.Sprintf(format, args...))
}

func InitDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	gormLogger := gormlogger.New(
		&gormLogWriter{l: logger.Log},
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Error,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
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
			logger.Log.Fatal().Err(err).Msg("Migration failed")
		}
	}

	DB = db
	logger.Log.Info().Msg("Database connected")
	return db, nil
}
