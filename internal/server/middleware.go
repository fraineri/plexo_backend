package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const UserIDKey key = "userID"

func (s *Server) authMiddleware(requiredPermission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"message":"Missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.settings.JWT.SecretKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"message":"Invalid token"}`, http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, `{"message":"Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			userID, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, `{"message":"Invalid user ID in token"}`, http.StatusUnauthorized)
				return
			}

			can, err := s.authService.Can(r.Context(), userID, requiredPermission)
			if err != nil {
				http.Error(w, `{"message":"Error checking permissions"}`, http.StatusInternalServerError)
				return
			}

			if !can {
				http.Error(w, `{"message":"Forbidden"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
