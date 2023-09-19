package service

import (
	"errors"
	"github.com/koteyye/brutalITSM-BE-Users/internal/models"
	"github.com/koteyye/brutalITSM-BE-Users/internal/postgres"
	"strings"
)

var AllDuplicate = errors.New("all duplicate")
var NoSettingsObject = errors.New("указанный объект отсутствует")

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
				Permissions: value.Permissions,
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

	r := 0
	for _, res := range result {
		if res.Success == false {
			r++
		} else {
			continue
		}
	}
	if r >= len(result) {
		return nil, AllDuplicate
	}

	return result, nil
}

func (s SettingsService) DeleteSettings(id []string, obj string) error {
	strings.ToLower(obj)
	var table string

	switch obj {
	case "job":
		table = "jobs"
		break
	case "org":
		table = "orgs"
		break
	case "role":
		table = "roles"
		break
	default:
		return NoSettingsObject
	}

	err := s.repo.DeleteSettings(id, table)
	return err
}

func (s SettingsService) EditSettings(set []models.Settings) ([]models.AddResult, error) {
	var result []models.AddResult

	var jobs []models.EditPq
	var orgs []models.EditPq
	var roles []models.RolesStr

	for _, value := range set {
		if value.SettingsObject == models.Job {
			jobs = append(jobs, models.EditPq{
				Id:   value.Id,
				Name: value.Name,
			})
		}
		if value.SettingsObject == models.Org {
			orgs = append(orgs, models.EditPq{
				Id:   value.Id,
				Name: value.Name,
			})
		}
		if value.SettingsObject == models.Role {
			roles = append(roles, models.RolesStr{
				Id:          value.Id,
				Name:        value.Name,
				Permissions: value.Permissions,
			})
		}
	}

	if jobs != nil {
		jobRes, err := s.repo.EditJob(jobs)
		if err != nil {
			return nil, err
		}
		result = append(result, jobRes...)
	}

	if orgs != nil {
		orgsRes, err := s.repo.EditOrg(orgs)
		if err != nil {
			return nil, err
		}
		result = append(result, orgsRes...)
	}

	if roles != nil {
		rolesRes, err := s.repo.EditRole(roles)
		if err != nil {
			return nil, err
		}
		result = append(result, rolesRes...)
	}
	return result, nil
}
