package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/sirupsen/logrus"
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
	row := u.db.QueryRow(query, user.Login, user.Password, user.Email, user.Lastname, user.Firstname, user.Surname, user.Job, user.Org, user.Role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (u userPostgres) CreateUserImg(userId string, avatar models.Avatar) (bool, error) {
	query := fmt.Sprintf("insert into user_img(user_id, mime_type, bucket_name, file_name)\nvalues ($1, $2, $3, $4);")
	row := u.db.QueryRow(query, userId, avatar.MimeType, avatar.BucketName, avatar.FileName)
	if err := row.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (u userPostgres) DeleteUserImg(imgId string) (bool, error) {
	query := fmt.Sprintf("update user_img set deleted_at = now() where id = $1;")
	row := u.db.QueryRow(query, imgId)
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

func (u userPostgres) GetUserAvatarS3(userId string) (models.SingleAvatars, error) {
	var s3Data models.SingleAvatars
	query := fmt.Sprintf("select id, bucket_name, file_name, mime_type from user_img where user_id = $1 and deleted_at is null")
	err := u.db.Get(&s3Data, query, userId)
	return s3Data, err
}

func (u userPostgres) CheckLogin(login string) (bool, error) {
	var duplicate bool
	query := fmt.Sprintf("select (case when (select \"login\" from users where \"login\" = $1 and deleted_at is null) is null then false else true end);")
	row := u.db.QueryRow(query, login)
	if err := row.Scan(&duplicate); err != nil {
		return false, err
	}
	return duplicate, nil
}

func (u userPostgres) GetUsers() ([]models.UserList, error) {
	var users []models.UserList

	query := fmt.Sprintf("select u.id, u.login, p.last_name, p.first_name, p.sur_name, j.name job_name, o.name org_name,\n       (select array_agg(r.name) from roles r join user_roles ur on r.id = ur.role_id where ur.user_id = u.id) role_list,\n       json_build_object('mimeType', ui.mime_type, 'bucketName', ui.bucket_name, 'fileName', ui.file_name) avatar\nfrom \"users\" u join persons p on u.id = p.user_id join jobs j on j.id = p.job_id join orgs o on o.id = p.org_id left join user_img ui on u.id = ui.user_id\nwhere u.deleted_at is null and ui.deleted_at is null;")
	err := u.db.Select(&users, query)

	return users, err
}

//func (u userPostgres) GetUserById(userId string) (models.UserList, error) {
//	var user models.UserList
//
//	query := fmt.Sprintf("select u.id,\n       u.login,\n       p.last_name,\n       p.first_name,\n       p.sur_name,\n  j.name job_name,\n o.name org_name,\n       (select array_agg(r.name)\n        from roles r\n                 join user_roles ur on r.id = ur.role_id\n        where ur.user_id = u.id)                                                                              role_list,\n       json_build_object('mimeType', ui.mime_type, 'bucketName', ui.bucket_name, 'fileName', ui.file_name) avatar\nfrom \"users\" u\n         join persons p on u.id = p.user_id\n join jobs j on j.id = p.job_id join orgs o on o.id = p.org_id         left join user_img ui on u.id = ui.user_id where u.id = $1 and u.deleted_at is null and ui.deleted_at is null;")
//	err := u.db.Get(&user, query, userId)
//
//	return user, err
//}

func (u userPostgres) GetRoles() ([]models.Roles, error) {
	var roles []models.Roles

	query := fmt.Sprintf("select id, name from roles")
	err := u.db.Select(&roles, query)

	return roles, err
}

func (u userPostgres) GetUserById(userId string) (models.UserList, error) {
	var user models.UserList

	query := sq.Select("*").From(fmt.Sprintf("getUserById('%v')", userId))
	sql, arg, err := query.ToSql()
	if err != nil {
		logrus.Fatalf("Какая-то дичь с запросом... %v", err)
	}
	logrus.Infof("%v", sql)
	logrus.Infof("arg: %v", arg)
	err1 := u.db.Get(&user, sql, userId)

	return user, err1
}
