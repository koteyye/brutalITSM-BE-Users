package service

import (
	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
	"errors"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"io"
)

type UserService struct {
	repo   repository.User
	s3repo *minio.Client
}

func NewUserService(repo repository.User, s3repo *minio.Client) *UserService {
	return &UserService{repo: repo, s3repo: s3repo}
}

func (u *UserService) CreateUser(user models.User) (string, error) {
	if user.Email == "" {
		user.Email = generateEmail(user.Login)
	}
	user.Password = generatePasswordHash(user.Password)
	return u.repo.CreateUser(user)
}

func (u *UserService) UploadFile(reader io.Reader, backetName string, fileName string) (string, error) {
	return "", nil
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

func (u *UserService) GetUserById(userId string) (models.UserList, error) {
	return u.repo.GetUserById(userId)
}
