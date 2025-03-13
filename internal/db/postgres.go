package db

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewConnection(dbURL string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dbURL), config)
	if err != nil {
		return nil, err
	}

	err = autoMigrate(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// autoMigrate automatically migrates the database schema
func autoMigrate(db *gorm.DB) error {
	log.Println("Running auto-migration")

	return db.AutoMigrate(
		&models.User{},
		&models.Author{},
		&models.Book{},
	)
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
