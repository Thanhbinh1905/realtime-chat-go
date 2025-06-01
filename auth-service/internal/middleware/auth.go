package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/config"
	"github.com/Thanhbinh1905/realtime-chat/auth-service/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

func (m *Middleware) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &utils.TokenClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
