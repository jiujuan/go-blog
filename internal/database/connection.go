package database

import (
	"go-blog/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the MySQL database
func Connect(databaseURL string) (*DB, error) {
	db, err := gorm.Open(mysql.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (these will be configurable via config)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return NewDB(db), nil
}

// Migrate runs database migrations using GORM AutoMigrate
func Migrate(db *DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.Article{},
		&models.Comment{},
		&models.Like{},
	)
}

// Close closes the database connection
func Close(db *DB) error {
	return db.Close()
}