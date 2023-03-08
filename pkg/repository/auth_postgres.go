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

func (r *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("select id from %s where login=$1 and password=$2", usersTable)
	err := r.db.Get(&user, query, login, password)

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

	query := fmt.Sprintf("select u.id,\n       u.login,\n       p.last_name,\n       p.first_name,\n       p.middle_name,\n       p.job_name,\n       p.org_name,\n       (select array_agg(r.name)\n        from roles r\n                 join user_roles ur on r.id = ur.role_id\n        where ur.user_id = u.id)                                                                              role_list,\n       json_build_object('mimeType', ui.mime_type, 'backetName', ui.backet_name, 'fileName', ui.file_name) avatar\nfrom \"user\" u\n         join person p on u.id = p.user_id\n         join user_img ui on u.id = ui.user_id\nwhere u.id = $1;")
	err := r.db.Get(&user, query, id)

	return user, err
}
