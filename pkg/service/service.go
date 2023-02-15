package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
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
