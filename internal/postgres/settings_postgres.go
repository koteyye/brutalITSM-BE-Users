package postgres

import (
	sql2 "database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	duplicateName = "duplicate name"
	deleteNull    = "not found role_id"
)

type rolePermissions struct {
	roleId string `db:"roles_id"`
}

type SettingsPostgres struct {
	db *sqlx.DB
}

func NewSettingsPostgres(db *sqlx.DB) *SettingsPostgres {
	return &SettingsPostgres{db: db}
}

func (s SettingsPostgres) AddOrg(orgNames pq.StringArray) ([]models.AddResult, error) {
	var orgRes []models.AddResult

	for _, que := range orgNames {
		var org models.AddResult

		query := sq.Insert("orgs").
			Columns("name").
			Values(que).
			Suffix("RETURNING \"id\", \"name\"").
			PlaceholderFormat(sq.Dollar)
		sql, arg, err := query.ToSql()
		if err != nil {
			logrus.Fatalf(toSqlErr, err)
		}
		row := s.db.QueryRow(sql, arg...)

		if err = row.Scan(&org.Id, &org.Name); err != nil {
			if checkDuplicateKey(err) == false {
				return nil, err
			}
			org.Name = que
			org.Success = false
			org.ErrorMessage = duplicateName
			orgRes = append(orgRes, org)
			continue
		}
		org.Success = true
		orgRes = append(orgRes, org)

	}

	return orgRes, nil
}

func (s SettingsPostgres) AddJob(jobNames pq.StringArray) ([]models.AddResult, error) {
	var jobRes []models.AddResult

	for _, que := range jobNames {
		var job models.AddResult

		query := sq.Insert("jobs").
			Columns("name").
			Values(que).
			Suffix("RETURNING \"id\", \"name\"").
			PlaceholderFormat(sq.Dollar)
		sql, arg, err := query.ToSql()
		if err != nil {
			logrus.Fatalf(toSqlErr, err)
		}
		row := s.db.QueryRow(sql, arg...)

		if err = row.Scan(&job.Id, &job.Name); err != nil {
			if checkDuplicateKey(err) == false {
				return nil, err
			}
			job.Name = que
			job.Success = false
			job.ErrorMessage = duplicateName
			jobRes = append(jobRes, job)
			continue
		}

		job.Success = true
		jobRes = append(jobRes, job)

	}

	return jobRes, nil
}

func (s SettingsPostgres) AddRole(roles []models.RolesStr) ([]models.AddResult, error) {
	var roleRes []models.AddResult
	for _, que := range roles {
		var role models.AddResult

		query := sq.Insert("roles").
			Columns("name").
			Values(que.Name).
			Suffix("RETURNING \"id\", \"name\"").
			PlaceholderFormat(sq.Dollar)

		sql, arg, err := query.ToSql()
		if err != nil {
			logrus.Fatalf(toSqlErr, err)
		}
		row := s.db.QueryRow(sql, arg...)

		if err = row.Scan(&role.Id, &role.Name); err != nil {
			if checkDuplicateKey(err) == false {
				return nil, err
			}
			role.Name = que.Name
			role.Success = false
			role.ErrorMessage = duplicateName
			roleRes = append(roleRes, role)
			continue
		}

		if que.Permissions != nil {
			for _, que2 := range que.Permissions {
				query := sq.Insert("role_permissions").
					Columns("roles_id, permission_id").Values(role.Id, que2).
					PlaceholderFormat(sq.Dollar)
				sql, arg, err := query.ToSql()
				if err != nil {
					logrus.Fatalf(toSqlErr, err)
				}
				row := s.db.QueryRow(sql, arg...)
				if err := row.Err(); err != nil {
					return nil, err
				}
			}
		}
		role.Success = true
		roleRes = append(roleRes, role)
	}
	return roleRes, nil
}

func checkDuplicateKey(err error) bool {
	var pgError *pq.Error
	if errors.As(err, &pgError) {
		if pgError.Code == "23505" {
			return true
		}
	}
	return false
}

func (s SettingsPostgres) DeleteSettings(id []string, table string) error {
	var deleteVar models.EditPq
	query := sq.Delete(table).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).Suffix("RETURNING id")
	sql, arg, err := query.ToSql()
	if err != nil {
		logrus.Fatalf(toSqlErr, err)
	}
	row := s.db.QueryRow(sql, arg...)
	return row.Scan(&deleteVar.Id)
}

func (s SettingsPostgres) EditJob(jobs []models.EditPq) ([]models.AddResult, error) {
	var jobRes []models.AddResult

	for _, que := range jobs {
		var job models.AddResult

		query := sq.Update("jobs").Set("name", que.Name).Where(sq.Eq{"id": que.Id}).PlaceholderFormat(sq.Dollar)
		sql, arg, err := query.ToSql()
		if err != nil {
			logrus.Fatalf(toSqlErr, err)
		}
		row := s.db.QueryRow(sql, arg...)
		if err := row.Err(); err != nil {
			if checkDuplicateKey(err) == false {
				return nil, err
			}
			job.Name = que.Name
			job.Success = false
			job.ErrorMessage = duplicateName
			jobRes = append(jobRes, job)
			continue
		}
		job.Name = que.Name
		job.Id = que.Id
		job.Success = true
		jobRes = append(jobRes, job)
	}
	return jobRes, nil
}

func (s SettingsPostgres) EditOrg(orgs []models.EditPq) ([]models.AddResult, error) {
	var orgRes []models.AddResult

	for _, que := range orgs {
		var org models.AddResult

		query := sq.Update("orgs").Set("name", que.Name).Where(sq.Eq{"id": que.Id}).PlaceholderFormat(sq.Dollar)
		sql, arg, err := query.ToSql()
		if err != nil {
			logrus.Fatalf(toSqlErr, err)
		}
		row := s.db.QueryRow(sql, arg...)
		if err := row.Err(); err != nil {
			if checkDuplicateKey(err) == false {
				return nil, err
			}
			org.Name = que.Name
			org.Success = false
			org.ErrorMessage = duplicateName
			orgRes = append(orgRes, org)
			continue
		}
		org.Name = que.Name
		org.Id = que.Id
		org.Success = true
		orgRes = append(orgRes, org)
	}
	return orgRes, nil
}

func (s SettingsPostgres) EditRole(roles []models.RolesStr) ([]models.AddResult, error) {
	var roleRes []models.AddResult

	for _, que := range roles {
		var role models.AddResult
		if que.Name != "" {
			query := sq.Update("roles").Set("name", que.Name).Where(sq.Eq{"id": que.Id}).PlaceholderFormat(sq.Dollar)
			sql, arg, err := query.ToSql()
			if err != nil {
				logrus.Fatalf(toSqlErr, err)
			}
			row := s.db.QueryRow(sql, arg...)
			if err := row.Err(); err != nil {
				if checkDuplicateKey(err) == false {
					return nil, err
				}
				role.Name = que.Name
				role.Success = false
				role.ErrorMessage = duplicateName
				roleRes = append(roleRes, role)
				continue
			}
			role.Name = que.Name
			role.Id = que.Id
			role.Success = true
			roleRes = append(roleRes, role)
		}
		if que.Permissions != nil {
			var roleId rolePermissions
			query := sq.Delete("role_permissions").Where(sq.Eq{"roles_id": que.Id}).Suffix("RETURNING \"roles_id\"").PlaceholderFormat(sq.Dollar)
			sql, arg, err := query.ToSql()
			if err != nil {
				logrus.Fatalf(toSqlErr, err)
			}
			row := s.db.QueryRow(sql, arg...)
			if err := row.Scan(&roleId.roleId); err != nil {
				switch err {
				case sql2.ErrNoRows:
					role.ErrorMessage = deleteNull
					break
				default:
					role.ErrorMessage = err.Error()
				}
				role.Name = que.Name
				role.Success = false
				roleRes = append(roleRes, role)
				continue
			}
			role.Success = true
			roleRes = append(roleRes, role)
		}
	}
	return roleRes, nil
}
