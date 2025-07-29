package database

import (
	"time"

	"go-blog/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectWithConfig establishes a connection to the MySQL database using config
func ConnectWithConfig(cfg *config.Config) (*DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.GetDatabaseURL()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool with config values
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings from config
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Second)

	return NewDB(db), nil
}