package database

import (
	"log"

	"github.com/obligate5818/bookrecs/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	log.Printf("Connecting to database with DSN: %q\n", dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate tables
	err = db.AutoMigrate(&models.Edition{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
