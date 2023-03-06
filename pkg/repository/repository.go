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
	DeleteUser(userId string) (bool, error)
	CheckLogin(login string) (bool, error)
	GetUsers() ([]models.UserList, error)
	GetUserById(userId string) (models.UserList, error)
}

type Repository struct {
	Authorization
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
	}
}
