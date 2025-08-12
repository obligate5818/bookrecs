package discord

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/obligate5818/bookrecs/internal/config"
	discordOauth "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

type DiscordTokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func NewOAuthConfig(cfg *config.Config) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  cfg.RedirectURI,
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Scopes:       []string{discordOauth.ScopeIdentify},
		Endpoint:     discordOauth.Endpoint,
	}
}

func OAuthStartHandler(oauthConf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := oauthConf.AuthCodeURL("state-string", oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
	}
}

func OAuthCallbackHandler(cfg *config.Config, oauthConf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Code not found in request", http.StatusBadRequest)
			return
		}

		// Exchange code for token
		token, err := oauthConf.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Use oauth2 client to get user info from Discord API
		client := oauthConf.Client(ctx, token)
		resp, err := client.Get("https://discord.com/api/users/@me")
		if err != nil {
			http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Discord API error: "+resp.Status, http.StatusInternalServerError)
			return
		}

		var user struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create JWT token with Discord user info
		jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":  user.ID,
			"name": user.Username,
			"exp":  time.Now().Add(24 * time.Hour).Unix(),
		})

		signedToken, err := jwtToken.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			http.Error(w, "Failed to sign JWT: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set JWT token in cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "Authorization",
			Value:    "Bearer " + signedToken,
			Path:     "/",
			HttpOnly: true,
			Secure:   false, // true if HTTPS
		})

		http.Redirect(w, r, "/", http.StatusFound)
	}
}
