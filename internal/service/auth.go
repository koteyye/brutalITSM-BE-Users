package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"os"
	"time"
)

const (
	tokenTTL = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

type AuthService struct {
	repo postgres.Authorization
}

func NewAuthService(repo postgres.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CheckRights(userId any, requierRole any) (bool, error) {

	role, err := s.repo.CheckRights(userId)
	if err != nil {
		return false, err
	}
	for _, i := range role {
		if i == requierRole || i == "admin" {
			return true, nil
		}
	}
	return false, errors.New("not enough rights")
}

func (s *AuthService) GenerateToken(login, password string) (string, error) {
	signKey := os.Getenv("SIGNING_KEY")
	auth, err := s.repo.Authentication(login, password)
	if err != nil {
		return "", errors.New("bad request")
	}

	if auth == false {
		return "", errors.New("invalid login or password")
	}

	user, err := s.repo.GetUser(login)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signKey))
}

func (s *AuthService) ParseToken(accessToken string) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		signKey := os.Getenv("SIGNING_KEY")
		return []byte(signKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return "", errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (s *AuthService) Me(id any) (models.UserList, error) {
	return s.repo.Me(id)
}

//func generatePasswordHash(password string) string {
//	hash := sha1.New()
//	hash.Write([]byte(password))
//
//	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
//}

func generateEmail(login string) string {
	mailDomainName := os.Getenv("DOMAIN_NAME")
	result := login + mailDomainName
	return result
}
