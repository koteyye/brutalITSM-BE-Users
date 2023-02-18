package repository

import (
	"brutalITSM-BE-Users/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CheckLogin(login string) (bool, error)
	CreateUser(user models.User) (string, error)
	GetUser(login, password string) (models.User, error)
}

type List interface {
}

type Repository struct {
	Authorization
	List
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}
