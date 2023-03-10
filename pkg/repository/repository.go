package repository

import (
	"brutalITSM-BE-Users/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CheckRights(userId any) ([]string, error)
	GetUser(login, password string) (models.User, error)
	Me(userId any) (models.UserList, error)
}

type User interface {
	CreateUser(user models.User) (string, error)
	CreateUserImg(userId string, user models.Avatar) (bool, error)
	DeleteUser(userId string) (bool, error)
	CheckLogin(login string) (bool, error)
	GetUsers() ([]models.UserList, error)
	GetUserById(userId string) (models.UserList, error)
	GetRoles() ([]string, error)
}

type Search interface {
	SearchJob(string) ([]string, error)
	SearchOrg(string) ([]string, error)
}

type Repository struct {
	Authorization
	User
	Search
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Search:        NewSearchPostgres(db),
	}
}
