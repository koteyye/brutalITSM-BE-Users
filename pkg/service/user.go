package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
	"errors"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) CreateUser(user models.User) (string, error) {
	if user.Email == "" {
		user.Email = generateEmail(user.Login)
	}

	user.Password = generatePasswordHash(user.Password)
	return u.repo.CreateUser(user)
}

func (u *UserService) CheckLogin(user models.User) (bool, error) {
	duplicate, err := u.repo.CheckLogin(user.Login)
	if err != nil {
		return false, err
	}
	logrus.Info(duplicate)
	if duplicate == true {
		return false, errors.New("duplicate login")
	}
	return true, nil
}

func (u *UserService) DeleteUser(userId string) (bool, error) {
	return u.repo.DeleteUser(userId)
}

func (u *UserService) GetUsers() ([]models.UserList, error) {
	return u.repo.GetUsers()
}
