package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"

	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
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
	return u.repo.CreateUser(user)
}

func (u *UserService) UploadFile(ctx context.Context, reader io.Reader, bucketName, fileName string, fileSize int64) (minio.UploadInfo, string, error) {
	info, err := u.s3repo.PutObject(ctx, bucketName, fileName, reader, fileSize, minio.PutObjectOptions{})

	if err != nil {
		return minio.UploadInfo{}, "", fmt.Errorf("cant upload file to s3")
	}

	mType, err := mimetype.DetectReader(reader)

	return info, mType.String(), nil
}

func (u *UserService) CreateUserImg(userId string, avatar models.Avatar) (bool, error) {
	return u.repo.CreateUserImg(userId, avatar)
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

func (u *UserService) GetRoles() ([]string, error) {
	return u.repo.GetRoles()
}
