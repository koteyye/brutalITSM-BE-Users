package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) Authentication(login, password string) (bool, error) {
	var ok bool
	query := fmt.Sprintf("select (password = crypt($1, password)) as pass from users where login = $2")
	row := r.db.QueryRow(query, password, login)
	if err := row.Scan(&ok); err != nil {
		return false, err
	}

	if ok == false {
		return false, nil
	}

	return true, nil
}

func (r *AuthPostgres) GetUser(login string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("select id from users where login = $1 and deleted_at is null")
	err := r.db.Get(&user, query, login)

	return user, err
}

func (r *AuthPostgres) CheckRights(userId any) ([]string, error) {
	var roleNames []string

	query := fmt.Sprintf("select get_user_roles($1)")
	err := r.db.Select(&roleNames, query, userId)

	return roleNames, err
}

func (r *AuthPostgres) Me(id any) (models.UserList, error) {
	var user models.UserList

	query := sq.Select("*").From("getUserById($1)").PlaceholderFormat(sq.Dollar)
	sql, _, err := query.ToSql()
	if err != nil {
		logrus.Fatalf("SQL query not builde %v", err)
	}
	err1 := r.db.Get(&user, sql, id)

	return user, err1
}
