package middleware

import (
	"context"
	"mpesa-finance/internal/auth"
	"net/http"
	"strings"
)

type contextKey string

const ClaimsKey contextKey = "claims"

//AuthMiddleware validates jwt tokens

func AuthMiddleware(authService *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Get authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Missing Authorization header" , "code": "UNAUTHORIZED"}`, http.StatusUnauthorized)
				return
			}
			//Extract token from header
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error": "Invalid Authorization header format" , "code": "UNAUTHORIZED"}`, http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]

			//Validate token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, `{"error":"Invalid token", "code":"UNAUTHORIZED"}`, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves claims from request context
func GetClaims(r *http.Request) (*auth.Claims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(*auth.Claims)
	return claims, ok
}
