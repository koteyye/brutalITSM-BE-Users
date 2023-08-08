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

func (s SettingsService) AddSettings(set []models.Settings) ([]models.AddResult, error) {
	var result []models.AddResult

	var jobs []string
	var orgs []string

	var roles []models.RolesStr

	for _, value := range set {
		if value.SettingsObject == models.Job {
			jobs = append(jobs, value.Name)
		}
		if value.SettingsObject == models.Org {
			orgs = append(orgs, value.Name)
		}
		if value.SettingsObject == models.Role {
			roles = append(roles, models.RolesStr{
				Name:        value.Name,
				Permissions: *value.Permissions,
			})
		}
	}

	if jobs != nil {
		jobRes, err := s.repo.AddJob(jobs)
		if err != nil {
			return nil, err
		}
		result = append(result, jobRes...)
	}

	if orgs != nil {
		orgRes, err := s.repo.AddOrg(orgs)
		if err != nil {
			return nil, err
		}
		result = append(result, orgRes...)
	}

	if roles != nil {
		roleRes, err := s.repo.AddRole(roles)
		if err != nil {
			return nil, err
		}
		result = append(result, roleRes...)
	}

	return result, nil
}

func (s SettingsService) DeleteSettings(id []string) (bool, error) {
	return s.DeleteSettings(id)
}

func (s SettingsService) EditSettings(setId string, set models.Settings) (string, error) {
	//TODO implement me
	panic("implement me")
}
