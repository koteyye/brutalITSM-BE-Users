package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
)

type Authorization interface {
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (string, error)
	CheckRights(userId any, requireRole any) (bool, error)
}

type User interface {
	CreateUser(user models.User) (string, error)
	DeleteUser(userId string) (bool, error)
	CheckLogin(user models.User) (bool, error)
	GetUsers() ([]models.UserList, error)
}

type Service struct {
	Authorization
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User),
	}
}
