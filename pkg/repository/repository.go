package repository

import (
	"brutalITSM-BE-Users/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
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
