package middleware

import (
	"context"
	"net/http"

	"github.com/slangeres/Vypaar/backend_API/internal/token"
)

type contextKey string

const emailKey contextKey = "email"

func AuthMiddleware(jwtMaker *token.JwtMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read JWT token from cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized: Missing token cookie", http.StatusUnauthorized)
				return
			}

			tokenStr := cookie.Value
			email, err := jwtMaker.VerifyToken(tokenStr)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), emailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
