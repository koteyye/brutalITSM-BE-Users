package repository

import (
	"brutalITSM-BE-Users/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

func (r *AuthPostgres) CheckRights(userId any) ([]string, error) {
	logrus.Info(userId)
	query := fmt.Sprintf("select get_user_roles($1)")
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var roleName []string
	logrus.Info(rows)
	for rows.Next() {
		if err := rows.Scan(&roleName); err != nil {
			return roleName, err
		}
	}
	if err = rows.Err(); err != nil {
		return roleName, err
	}

	logrus.Info(roleName)

	return roleName, nil
}
