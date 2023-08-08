package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type settingsPostgres struct {
	db *sqlx.DB
}

func NewSettingsPostgres(db *sqlx.DB) *settingsPostgres {
	return &settingsPostgres{db: db}
}

func (s settingsPostgres) AddOrg(orgNames pq.StringArray) ([]models.AddResult, error) {
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
			return nil, err
		}

		orgRes = append(orgRes, org)

	}

	return orgRes, nil
}

func (s settingsPostgres) AddJob(jobNames pq.StringArray) ([]models.AddResult, error) {
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
			return nil, err
		}

		jobRes = append(jobRes, job)

	}

	return jobRes, nil
}

func (s settingsPostgres) AddRole(roles []models.RolesStr) ([]models.AddResult, error) {
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
			return nil, err
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
		roleRes = append(roleRes, role)
	}
	return roleRes, nil
}

func (s settingsPostgres) DeleteSettings(id []string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s settingsPostgres) EditSettings(id []string, set models.Settings) (bool, error) {
	//TODO implement me
	panic("implement me")
}
