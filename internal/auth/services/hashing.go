package services

import "golang.org/x/crypto/bcrypt"

type HashingService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type bcryptHashingService struct{}

func NewHashingService() HashingService {
	return &bcryptHashingService{}
}

func (s *bcryptHashingService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *bcryptHashingService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
