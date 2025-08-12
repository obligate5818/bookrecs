package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/obligate5818/bookrecs/internal/models"
	"github.com/obligate5818/bookrecs/internal/openlibrary"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	userID, ok1 := r.Context().Value("discord_user_id").(string)
	username, ok2 := r.Context().Value("discord_username").(string)

	if !ok1 || !ok2 || userID == "" || username == "" {
		http.Error(w, "User information not found in context", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head><title>Home</title></head>
		<body>
			<h1>Hello, %s!</h1>
			<p>Your Discord ID is %s.</p>
			<p><a href="/books">Go to books list</a></p>
			<p><a href="/isbn">Go to POST ISBN form</a></p>
		</body>
		</html>
	`, username, userID)
}

func GetIsbnForm(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head><title>Enter ISBN</title></head>
	<body>
		<h1>Enter ISBN</h1>
		<form method="POST" action="/isbn">
			<label for="isbn">ISBN:</label>
			<input type="text" id="isbn" name="isbn" required>
			<button type="submit">Submit</button>
		</form>
	</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

type FetchEditionFunc func(ctx context.Context, isbn string) (*openlibrary.Edition, error)

func GetBooks(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editions []models.Edition
		if err := db.Find(&editions).Error; err != nil {
			http.Error(w, "failed to fetch books", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(editions); err != nil {
			http.Error(w, "failed to encode books", http.StatusInternalServerError)
		}
	}
}

func PostIsbn(db *gorm.DB, fetchEdition FetchEditionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed to parse form", http.StatusBadRequest)
			return
		}

		isbn := strings.TrimSpace(r.FormValue("isbn"))
		if isbn == "" {
			http.Error(w, "missing isbn", http.StatusBadRequest)
			return
		}

		olEdition, err := fetchEdition(r.Context(), isbn)
		if err != nil {
			http.Error(w, "fetch failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		edition := olEdition.ToInternalModel()

		err = db.Where(models.Edition{Key: edition.Key}).FirstOrCreate(edition).Error
		if err != nil {
			http.Error(w, "store failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, edition.Key, http.StatusSeeOther)
	}
}

func GetEdition(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		keyID := vars["key__id"]
		key := "/books/" + keyID

		ed, err := getEditionByKey(db, key)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(ed)
	}
}

func getEditionByKey(db *gorm.DB, key string) (*models.Edition, error) {
	var ed models.Edition
	if err := db.Where("key = ?", key).First(&ed).Error; err != nil {
		return nil, err
	}
	return &ed, nil
}
