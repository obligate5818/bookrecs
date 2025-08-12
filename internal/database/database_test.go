package database_test

import (
	"testing"

	"github.com/obligate5818/bookrecs/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestSQLiteMigrateFetchStoreQuery(t *testing.T) {
	// 1. Open in-memory SQLite DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // or logger.Silent, logger.Error, logger.Warn
	})
	if err != nil {
		t.Fatal(err)
	}

	// 2. Run AutoMigrate on the Edition model
	if err := db.AutoMigrate(&models.Edition{}); err != nil {
		t.Fatal(err)
	}

}
