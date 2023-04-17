package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
)

type searchPostgres struct {
	db *sqlx.DB
}

func NewSearchPostgres(db *sqlx.DB) *searchPostgres {
	return &searchPostgres{db: db}
}

func (s searchPostgres) SearchJob(searchQuery string) ([]models.SearchResult, error) {
	var searchResult []models.SearchResult
	query := fmt.Sprintf("select id, name from jobs where lower(name) like lower($1);")
	err := s.db.Select(&searchResult, query, "%"+searchQuery+"%")

	return searchResult, err
}

func (s searchPostgres) SearchOrg(searchQuery string) ([]models.SearchResult, error) {
	var searchResult []models.SearchResult
	query := fmt.Sprintf("select id, name from orgs where lower(name) like lower($1);")
	err := s.db.Select(&searchResult, query, "%"+searchQuery+"%")

	return searchResult, err
}
