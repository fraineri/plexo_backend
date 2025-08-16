package services

import (
	"fmt"
	"time"

	"github.com/fraineri/plexo_backend/internal/auth/entities"
	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(user *entities.User, roles []string) (string, error)
}

type jwtService struct {
	secretKey    string
	session_time int
}

func NewJWTService(secretKey string, session_time int) JWTService {
	return &jwtService{
		secretKey:    secretKey,
		session_time: session_time,
	}
}

func (s *jwtService) GenerateToken(user *entities.User, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"roles": roles,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * time.Duration(s.session_time)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
