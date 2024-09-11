package database

import (
	"go-blog/internal/domain/models"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BlogPost = models.BlogPost

func NewPostgresConnection(connStr string) (*gorm.DB, error) {
	DB := &gorm.DB{}
	// set variables
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = DB.AutoMigrate(&BlogPost{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return DB, nil
}

func CloseDatabase(DB *gorm.DB) {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting underlying SQL DB:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Error closing database connection:", err)
	}

	log.Println("Database and listener closed successfully")
}
