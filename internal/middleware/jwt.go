package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/learies/gofermart/internal/services"
)

var jwtService = services.NewJWTService()

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		cookie, err := r.Cookie("token")
		if err == nil {
			tokenString = cookie.Value
		}

		if tokenString != "" {
			userID, err := jwtService.VerifyToken(tokenString)
			if err != nil {
				log.Println("Invalid token")
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}