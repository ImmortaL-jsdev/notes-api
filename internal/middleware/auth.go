package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/ImmortaL-jsdev/notes-api/internal/auth"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte(os.Getenv("JWT_SECRET"))

		if len(secret) == 0 {
			secret = []byte("supersecret")
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if tokenString == authHeader {
			http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(tokenString, secret)

		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		userID, ok := claims["sub"].(string)

		if !ok {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
