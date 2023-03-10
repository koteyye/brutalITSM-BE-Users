package service

import (
	"context"
	"io"

	"brutalITSM-BE-Users/models"
	"brutalITSM-BE-Users/pkg/repository"
	"github.com/minio/minio-go/v7"
)

type Authorization interface {
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (string, error)
	CheckRights(userId any, requireRole any) (bool, error)
	Me(id any) (models.UserList, error)
}

type User interface {
	CreateUser(user models.User) (string, error)
	CreateUserImg(userId string, avatar models.Avatar) (bool, error)
	DeleteUser(userId string) (bool, error)
	CheckLogin(user models.User) (bool, error)
	GetUsers() ([]models.UserList, error)
	GetUserById(userId string) (models.UserList, error)
	UploadFile(ctx context.Context, reader io.Reader, backetName string, filename string, fileSize int64) (minio.UploadInfo, string, error)
	GetRoles() ([]string, error)
}

type Search interface {
	SearchJob(string) ([]string, error)
	SearchOrg(string) ([]string, error)
}

type Service struct {
	Authorization
	User
	Search
}

func NewService(repos *repository.Repository, s3repo *minio.Client) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User, s3repo),
		Search:        NewSearchService(repos.Search),
	}
}
