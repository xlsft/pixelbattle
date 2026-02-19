package database

import (
	"fmt"
	"os"

	"github.com/xlsft/pixelbattle/src/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitializeDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	fmt.Println("Connecting to database with DSN:", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	Db = db
	db.AutoMigrate(&models.PixelModel{})
	return db, nil
}

func UseDb() *gorm.DB {
	return Db
}
