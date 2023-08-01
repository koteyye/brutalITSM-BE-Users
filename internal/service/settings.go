package service

import (
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
)

type SettingsService struct {
	repo postgres.Settings
}

func NewSettingsService(repo postgres.Settings) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s SettingsService) AddSettings(set models.Settings) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (s SettingsService) DeleteSettings(id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (s SettingsService) EditSettings(setId string, set models.Settings) (string, error) {
	//TODO implement me
	panic("implement me")
}
