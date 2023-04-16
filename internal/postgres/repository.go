package postgres

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type User interface {
}

type Search interface {
}

type Repository struct {
	Authorization
	User
	Search
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
