package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/obligate5818/bookrecs/internal/config"
	"github.com/obligate5818/bookrecs/internal/database"
	"github.com/obligate5818/bookrecs/internal/discord"
	"github.com/obligate5818/bookrecs/internal/handlers"
	"github.com/obligate5818/bookrecs/internal/openlibrary"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	oauth_conf := discord.NewOAuthConfig(cfg)

	r := mux.NewRouter()

	// OAuth routes
	r.HandleFunc("/oauth/start", discord.OAuthStartHandler(oauth_conf))
	r.HandleFunc("/oauth/callback", discord.OAuthCallbackHandler(cfg, oauth_conf))

	r.HandleFunc("/", handlers.AuthMiddleware(cfg.JWTSecret, handlers.GetHome))

	r.HandleFunc("/isbn", handlers.AuthMiddleware(cfg.JWTSecret, handlers.GetIsbnForm)).Methods("GET")
	r.HandleFunc("/isbn", handlers.AuthMiddleware(cfg.JWTSecret, handlers.PostIsbn(db, openlibrary.FetchEdition))).Methods("POST")
	r.HandleFunc("/books", handlers.AuthMiddleware(cfg.JWTSecret, handlers.GetBooks(db))).Methods("GET")
	r.HandleFunc("/books/{key__id}", handlers.GetEdition(db)).Methods("GET")

	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}).Methods("GET")

	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

// API routes with auth
// r.HandleFunc("/isbn", handlers.AuthMiddleware(cfg.JWTSecret, handlers.PostIsbn(db))).Methods("POST")
// r.HandleFunc("/books/{key__id}.json", handlers.GetEdition(db)).Methods("GET")
