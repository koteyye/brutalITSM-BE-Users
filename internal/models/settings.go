package models

import "github.com/lib/pq"

const (
	Role string = "role"
	Job  string = "job"
	Org  string = "org"
)

type Settings struct {
	Id             string         `json:"id" db:"id"`
	Name           string         `json:"name" db:"name"`
	SettingsObject string         `json:"settingsObject" db:"settingsObject"`
	Permissions    pq.StringArray `json:"permissions" db:"permissions"`
}

type RolesStr struct {
	Id          string
	Name        string
	Permissions pq.StringArray
}

type AddResult struct {
	Name         string `json:"name" db:"name"`
	Id           string `json:"id" db:"id"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

type EditPq struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}
