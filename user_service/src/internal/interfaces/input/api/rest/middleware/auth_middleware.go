package middleware

import (
	"context"
	"net/http"
	"strings"
	"user_service/pkg/generatejwt"
)

type contextKey string

const UserIDContextKey = contextKey("userID")

// middleware for JWT authentication.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// header expected in the format "Bearer <token>".
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			tokenString := headerParts[1]

			// validating the token
			claims, err := generatejwt.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// if it reaches here, token is valid...hence, adding UserID to request context
			ctx := context.WithValue(r.Context(), UserIDContextKey, claims.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
