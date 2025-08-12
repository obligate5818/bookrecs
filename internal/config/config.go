package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	JWTSecret    string
	Port         string
	SessionKey   string
	SafeGuildIDs []string
	DatabaseURL  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	must := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("missing env var %s", k)
		}
		return v
	}

	return &Config{
		ClientID:     must("DISCORD_CLIENT_ID"),
		ClientSecret: must("DISCORD_CLIENT_SECRET"),
		RedirectURI:  must("DISCORD_REDIRECT_URI"),
		JWTSecret:    must("BOOKRECS_JWT_SECRET"),
		Port:         must("BOOKRECS_PORT"),
		SessionKey:   must("BOOKRECS_SESSION_KEY"),
		SafeGuildIDs: strings.Split(must("BOOKRECS_SAFE_GUILD_ID"), ","),
		DatabaseURL:  must("BOOKRECS_DATABASE_URL"),
	}
}
