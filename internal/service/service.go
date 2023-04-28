package service

import (
	"context"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"github.com/minio/minio-go/v7"
	"io"
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
	UpdateUserImg(ctx context.Context, userId string, avatar models.Avatar) (bool, error)
	DeleteUser(ctx context.Context, userId string) (bool, error)
	CheckLogin(user models.User) (bool, error)
	GetUsers() ([]models.UserList, error)
	GetUserById(userId string) (models.UserList, error)
	UploadFile(ctx context.Context, reader io.Reader, bucketName string, filename string, fileSize int64) (minio.UploadInfo, string, error)
	GetRoles() ([]models.Roles, error)
	GetUserList(userId []string) ([]models.UserShortList, error)
}

type Search interface {
	SearchJob(string) ([]models.SearchResult, error)
	SearchOrg(string) ([]models.SearchResult, error)
}

type Service struct {
	Authorization
	User
	Search
}

func NewService(repos *postgres.Repository, s3 *minio.Client) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		User:          NewUserService(repos.User, s3),
		Search:        NewSearchService(repos.Search),
	}
}
