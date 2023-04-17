package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
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

	query := fmt.Sprintf("select u.id, u.login, p.last_name, p.first_name, j.name job_name, o.name org_name, (select array_agg(r.name) from roles r join user_roles ur on r.id = ur.role_id where ur.user_id = u.id) role_list, json_build_object('mimeType', ui.mime_type, 'bucketName', ui.bucket_name, 'fileName', ui.file_name) avatar from users u join persons p on u.id = p.user_id left join user_img ui on u.id = ui.user_id join jobs j on p.job_id = j.id join orgs o on p.org_id = o.id where u.id = $1 and ui.deleted_at is null;")
	err := r.db.Get(&user, query, id)

	return user, err
}
