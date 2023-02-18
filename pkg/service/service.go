package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (string, error)
	CheckLogin(user models.User) (bool, error)
}

type List interface {
}

type Service struct {
	Authorization
	List
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
