package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
)

type settingsPostgres struct {
	db *sqlx.DB
}

func NewSettingsPostgres(db *sqlx.DB) *settingsPostgres {
	return &settingsPostgres{db: db}
}

func (s settingsPostgres) AddSettings(set models.Settings) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s settingsPostgres) DeleteSettings(id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s settingsPostgres) EditSettings(id string, set models.Settings) (bool, error) {
	//TODO implement me
	panic("implement me")
}
