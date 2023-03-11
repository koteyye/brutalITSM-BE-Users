package repository

import (
	"brutalITSM-BE-Users/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type userPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *userPostgres {
	return &userPostgres{db: db}
}

func (u userPostgres) CreateUser(user models.User) (string, error) {
	var id string
	query := fmt.Sprintf("select createUser($1, $2, $3, $4, $5, $6, $7, $8, $9);")
	row := u.db.QueryRow(query, user.Login, user.Password, user.Email, user.Lastname, user.Firstname, user.Middlename, user.Jobname, user.Orgname, user.RoleName)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (u userPostgres) CreateUserImg(userId string, avatar models.Avatar) (bool, error) {
	query := fmt.Sprintf("insert into user_img(user_id, mime_type, backet_name, file_name)\nvalues ($1, $2, $3, $4);")
	row := u.db.QueryRow(query, userId, avatar.MimeType, avatar.BacketName, avatar.FileName)
	if err := row.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (u userPostgres) DeleteUser(userId string) (bool, error) {
	var ok bool
	query := fmt.Sprintf("select delete_user($1);")
	row := u.db.QueryRow(query, userId)
	if err := row.Scan(&ok); err != nil {
		return false, err
	}
	return ok, nil
}

func (u userPostgres) CheckLogin(login string) (bool, error) {
	var duplicate bool
	query := fmt.Sprintf("select (case when (select \"login\" from %s where \"login\" = $1) is null then false else true end\n           );", usersTable)
	row := u.db.QueryRow(query, login)
	if err := row.Scan(&duplicate); err != nil {
		return false, err
	}
	return duplicate, nil
}

func (u userPostgres) GetUsers() ([]models.UserList, error) {
	var users []models.UserList

	query := fmt.Sprintf("select u.id,\n       u.login,\n       p.last_name,\n       p.first_name,\n       p.middle_name,\n       p.job_name,\n       p.org_name,\n       (select array_agg(r.name)\n        from roles r\n                 join user_roles ur on r.id = ur.role_id\n        where ur.user_id = u.id)                                                                              role_list,\n       json_build_object('mimeType', ui.mime_type, 'backetName', ui.backet_name, 'fileName', ui.file_name) avatar\nfrom \"user\" u\n         join person p on u.id = p.user_id\n         left join user_img ui on u.id = ui.user_id;")
	err := u.db.Select(&users, query)

	return users, err
}

func (u userPostgres) GetUserById(userId string) (models.UserList, error) {
	var user models.UserList

	query := fmt.Sprintf("select u.id, u.login, last_name, first_name, middle_name, job_name, org_name,\n       (select array_agg(r.name) from roles r join user_roles ur on r.id = ur.role_id where ur.user_id = u.id) role_list\nfrom \"user\" u\njoin person p on p.user_id = u.id\nwhere u.id = $1;")
	err := u.db.Get(&user, query, userId)

	return user, err
}
