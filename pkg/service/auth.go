package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
	"crypto/sha1"
	"fmt"
)

const salt = "io9ghtreh5erhefsdgewrggdfc"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user models.User) (string, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
