package repository

import (
	"brutalITSM-BE-Users/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CheckLogin(login string) (bool, error) {
	var duplicate bool
	query := fmt.Sprintf("select (case when (select \"login\" from %s where \"login\" = $1) is null then false else true end\n           );", usersTable)
	err := r.db.Get(&duplicate, query, login)
	return duplicate, err
}

func (r *AuthPostgres) CreateUser(user models.User) (string, error) {
	var id string
	query := fmt.Sprintf("select createUser($1, $2, $3, $4, $5, $6, $7, $8, $9);")
	row := r.db.QueryRow(query, user.Login, user.Password, user.Email, user.Lastname, user.Firstname, user.Middlename, user.Jobname, user.Orgname, user.RoleName)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("select id from %s where login=$1 and password=$2", usersTable)
	err := r.db.Get(&user, query, login, password)

	return user, err
}
