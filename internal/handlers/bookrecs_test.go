package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/obligate5818/bookrecs/internal/models"
	"github.com/obligate5818/bookrecs/internal/openlibrary"
)

func TestPostIsbn(t *testing.T) {
	// Setup in-memory sqlite DB or test DB for GORM
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // or logger.Silent, logger.Error, logger.Warn
	})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&models.Edition{})
	if err != nil {
		t.Fatal(err)
	}

	// Mock fetch function reads from local testdata file
	mockFetch := func(ctx context.Context, isbn string) (*openlibrary.Edition, error) {
		data, err := os.ReadFile("../testdata/edition.json")
		if err != nil {
			return nil, err
		}
		var ed openlibrary.Edition
		if err := json.Unmarshal(data, &ed); err != nil {
			return nil, err
		}
		return &ed, nil
	}

	handler := PostIsbn(db, mockFetch)

	form := url.Values{}
	form.Set("isbn", "1234") // isbn does not matter for this test
	req := httptest.NewRequest("POST", "/isbn", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status 303 See Other, got %d", resp.StatusCode)
	}

	location := resp.Header.Get("Location")
	if location == "" {
		t.Error("expected redirect location header")
	}

	// Optionally verify the DB has the stored edition
	var count int64
	db.Model(&models.Edition{}).Count(&count)
	if count != 1 {
		t.Errorf("expected 1 edition in DB, got %d", count)
	}
}
