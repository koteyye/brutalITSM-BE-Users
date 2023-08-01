package models

import "github.com/lib/pq"

type Settings struct {
	Id             string          `json:"id" db:"id"`
	Name           string          `json:"name" db:"name"`
	SettingsObject string          `json:"settingsObject" db:"settingsObject"`
	Permissions    *pq.StringArray `json:"permissions" db:"permissions"`
}
