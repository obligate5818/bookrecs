package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 {
				tokenStr = parts[1]
			}
		} else {
			// Try from cookie
			cookie, err := r.Cookie("Authorization")
			if err == nil {
				tokenStr = strings.TrimPrefix(cookie.Value, "Bearer ")
			}
		}

		if tokenStr == "" {
			http.Redirect(w, r, "/oauth/start", http.StatusFound)
			return
		}

		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "discord_user_id", claims["sub"])
		ctx = context.WithValue(ctx, "discord_username", claims["name"])
		next(w, r.WithContext(ctx))
	}
}
