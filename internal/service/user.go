package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"io"

	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/minio/minio-go/v7"
)

type UserService struct {
	repo   postgres.User
	s3repo *minio.Client
}

func NewUserService(repo postgres.User, s3repo *minio.Client) *UserService {
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

func (u *UserService) UpdateUserImg(ctx context.Context, userId string, avatar models.Avatar) (bool, error) {
	s3file, err := u.repo.GetUserAvatarS3(userId)
	switch err {
	case nil:
		break
	case sql.ErrNoRows:
		return u.repo.CreateUserImg(userId, avatar)
	default:
		return false, err
	}

	err1 := u.s3repo.RemoveObject(ctx, s3file.BucketName, s3file.FileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err1 != nil {
		return false, err1
	}

	_, err2 := u.repo.DeleteUserImg(s3file.ImgId)
	if err2 != nil {
		return false, err2
	}

	return u.repo.CreateUserImg(userId, avatar)
}

func (u *UserService) CheckLogin(user models.User) (bool, error) {
	duplicate, err := u.repo.CheckLogin(user.Login)
	if err != nil {
		return false, err
	}
	if duplicate == true {
		return false, errors.New("duplicate login")
	}
	return true, nil
}

func (u *UserService) DeleteUser(ctx context.Context, userId string) (bool, error) {
	s3file, err := u.repo.GetUserAvatarS3(userId)

	switch err {
	case nil:
		err2 := u.s3repo.RemoveObject(ctx, s3file.BucketName, s3file.FileName, minio.RemoveObjectOptions{ForceDelete: true})
		if err2 != nil {
			return false, err2
		}
		break
	case sql.ErrNoRows:
		break
	default:
		return false, err
	}

	return u.repo.DeleteUser(userId)
}

func (u *UserService) GetUsers() ([]models.UserList, error) {
	return u.repo.GetUsers()
}

func (u *UserService) GetUserById(userId string) (models.UserList, error) {
	return u.repo.GetUserById(userId)
}

func (u *UserService) GetRoles() ([]models.Roles, error) {
	return u.repo.GetRoles()
}
