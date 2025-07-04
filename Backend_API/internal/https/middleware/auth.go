package middleware

import (
	"context"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/slangeres/Vypaar/backend_API/internal/token"
)

type contextKey string

const userClaimsKey contextKey = "userClaims"

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

			var claims token.VypaarClaim

			jwttoken, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtMaker.SecretKey), nil
			})
			if err != nil || !jwttoken.Valid {
				http.Error(w, "Unauthorised Access", http.StatusUnauthorized)
				return

			}
			ctx := context.WithValue(r.Context(), userClaimsKey, &claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
func UserClaimsFromContext(ctx context.Context) (*token.VypaarClaim, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*token.VypaarClaim)
	return claims, ok
}
