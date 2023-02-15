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

func (r *AuthPostgres) CreateUser(user models.User) (string, error) {
	var id string
	var personId string
	query := fmt.Sprintf("INSERT INTO %s (login, password) VALUES ($1, $2) returning id", usersTable)
	query2 := fmt.Sprintf("INSERT INTO %s (last_name, first_name, middle_name, job_name, org_name, user_id) VALUES ($1, $2, $3, $4, $5, $6) returning id", personTable)

	row := r.db.QueryRow(query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return "", err
	} else {
		row := r.db.QueryRow(query2, user.Lastname, user.Firstname, user.Middlename, user.Jobname, user.Orgname, id)
		if err := row.Scan(&personId); err != nil {
			return "", err
		}
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("select id from %s where login=$1 and password=$2", usersTable)
	err := r.db.Get(&user, query, login, password)

	return user, err
}
