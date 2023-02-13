package user

import (
	"brutalITSM-BE-Users/internal/user"
	"brutalITSM-BE-Users/pkg/client/postgresql"
	"brutalITSM-BE-Users/pkg/logging"
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, user *user.User) error {

	q := `
	INSERT INTO
	    "user" (login, password) 
	VALUES 
	    ($1, crypt($2, gen_salt('md5'))) 
	RETURNING id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	if err := r.client.QueryRow(ctx, q, user.Login, user.Password).Scan(&user.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLStater: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil
		}
		r.logger.Info(err)
		return err
	}

	return nil
}

func (r *repository) CreatePerson(ctx context.Context, person *user.Person) error {
	q := `
	insert into
		person (last_name, first_name, middle_name, job_name, org_name, user_id)
	values
		($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", q))
	if err := r.client.QueryRow(ctx, q, person.LastName, person.FirstName, person.MiddleName, person.JobName, person.OrgName, person.UserId).Scan(&person.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLStater: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil
		}

		return err
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
