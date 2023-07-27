package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/sirupsen/logrus"
)

const toSqlErr = "SQL query not build %v"

type userPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *userPostgres {
	return &userPostgres{db: db}
}

func (u userPostgres) CreateUser(user models.User) (string, error) {
	var id string

	query := sq.Select("createUser($1, $2, $3, $4, $5, $6, $7, $8, $9)").PlaceholderFormat(sq.Dollar)
	sql, _, err := query.ToSql()
	if err != nil {
		logrus.Fatalf(toSqlErr, err)
	}
	row := u.db.QueryRow(sql, user.Login, user.Password, user.Email, user.Lastname, user.Firstname, user.Surname, user.Job, user.Org, user.Role)
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

//var users []models.UserShortList
//
//query := sq.Select("p.user_id, p.last_name, p.first_name, p.sur_name,  ui.mime_type, ui.bucket_name, ui.file_name").
//From("persons p").LeftJoin("user_img ui on p.user_id = ui.user_id").
//Where(sq.Eq{"p.user_id": usersId}).
//PlaceholderFormat(sq.Dollar)
//sql, args, err := query.ToSql()
//if err != nil {
//logrus.Fatalf("SQL query not builde %v", err)
//}
//err1 := u.db.Select(&users, sql, args...)
//return users, err1

func (u userPostgres) CreateUserImg(userId string, avatar models.Avatar) (bool, error) {
	query := sq.Insert("user_img").
		Columns("user_id, mime_type, bucket_name, file_name").
		Values(userId, avatar.MimeType, avatar.BucketName, avatar.FileName).
		PlaceholderFormat(sq.Dollar)
	sql, arg, err := query.ToSql()
	if err != nil {
		logrus.Fatalf(toSqlErr, err)
	}
	row := u.db.QueryRow(sql, arg...)
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

	query := sq.Select("*").From("getUserList()")
	sql, _, err := query.ToSql()
	if err != nil {
		logrus.Fatalf("SQL query not builde %v", err)
	}
	err1 := u.db.Select(&users, sql)
	return users, err1

}

func (u userPostgres) GetRoles() ([]models.Roles, error) {
	var roles []models.Roles

	query := fmt.Sprintf("select id, name from roles")
	err := u.db.Select(&roles, query)

	return roles, err
}

func (u userPostgres) GetUserById(userId string) (models.UserList, error) {
	var user models.UserList
	query := sq.Select("*").From("getUserById($1)").PlaceholderFormat(sq.Dollar)
	sql, _, err := query.ToSql()
	if err != nil {
		logrus.Fatalf("SQL query not builde %v", err)
	}
	err1 := u.db.Get(&user, sql, userId)

	return user, err1
}

func (u userPostgres) GetUserList(usersId []string) ([]models.UserShortList, error) {
	var users []models.UserShortList

	query := sq.Select("p.user_id, p.last_name, p.first_name, p.sur_name,  ui.mime_type, ui.bucket_name, ui.file_name").
		From("persons p").LeftJoin("user_img ui on p.user_id = ui.user_id").
		Where(sq.Eq{"p.user_id": usersId, "ui.deleted_at": nil}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		logrus.Fatalf("SQL query not builde %v", err)
	}
	err1 := u.db.Select(&users, sql, args...)
	return users, err1
}
