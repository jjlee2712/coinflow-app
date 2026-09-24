package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/jjlee2712/coinflow/backend/service"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(jwtSecret string) func(http.Handler)  http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer "){
        http.Error(w, "missing or invalid token", http.StatusUnauthorized)
        return
			}
      tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

      userID, err := service.ValidateToken(tokenStr, jwtSecret)
      if err != nil {
        http.Error(w, "invalid token", http.StatusUnauthorized)
        return
      }
      ctx := context.WithValue(r.Context(), UserIDKey, userID)
      next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}