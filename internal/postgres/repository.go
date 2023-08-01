package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
)

type Authorization interface {
	CheckRights(userId any) ([]string, error)
	GetUser(login string) (models.User, error)
	Authentication(login, password string) (bool, error)
	Me(userId any) (models.UserList, error)
}

type User interface {
	CreateUser(user models.User) (string, error)
	CreateUserImg(userId string, user models.Avatar) (bool, error)
	DeleteUserImg(imgId string) (bool, error)
	DeleteUser(userId string) (bool, error)
	CheckLogin(login string) (bool, error)
	GetUsers() ([]models.UserList, error)
	GetUserById(userId string) (models.UserList, error)
	GetRoles() ([]models.Roles, error)
	GetUserAvatarS3(userId string) (models.SingleAvatars, error)
	GetUserList(usersId []string) ([]models.UserShortList, error)
}

type Search interface {
	SearchJob(string) ([]models.SearchResult, error)
	SearchOrg(string) ([]models.SearchResult, error)
}

type Settings interface {
	AddSettings(set models.Settings) ([]string, error)
	DeleteSettings(id string) (bool, error)
	EditSettings(id string, set models.Settings) (bool, error)
}

type Repository struct {
	Authorization
	User
	Search
	Settings
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Search:        NewSearchPostgres(db),
		Settings:      NewSettingsPostgres(db),
	}
}
