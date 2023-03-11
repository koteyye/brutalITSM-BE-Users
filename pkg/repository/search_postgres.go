package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type searchPostgres struct {
	db *sqlx.DB
}

func NewSearchPostgres(db *sqlx.DB) *searchPostgres {
	return &searchPostgres{db: db}
}

func (s searchPostgres) SearchJob(searchQuery string) ([]string, error) {
	var searchResult []string
	query := fmt.Sprintf("select job_name from person where lower(job_name) like lower($1) group by job_name;")
	err := s.db.Select(&searchResult, query, "%"+searchQuery+"%")

	return searchResult, err
}

func (s searchPostgres) SearchOrg(searchQuery string) ([]string, error) {
	var searchResult []string
	query := fmt.Sprintf("select org_name from person where lower(org_name) like lower($1) group by org_name;")
	err := s.db.Select(&searchResult, query, "%"+searchQuery+"%")

	return searchResult, err
}
