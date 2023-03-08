package models

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	Id         string `json:"-" db:"id"`
	Login      string `json:"login" binding:"required"`
	Email      string `json:"email"`
	Password   string `json:"password" binding:"required"`
	Lastname   string `json:"lastname"`
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Jobname    string `json:"jobname"`
	Orgname    string `json:"orgname"`
	Userid     string `json:"-"`
	RoleName   string `json:"roleName" db:"name"`
}

type UserList struct {
	Id         string         `json:"id" db:"id"`
	Login      string         `json:"login" db:"login"`
	Lastname   string         `json:"lastname" db:"last_name"`
	Firstname  string         `json:"firstname" db:"first_name"`
	Middlename string         `json:"middlename" db:"middle_name"`
	Jobname    string         `json:"jobname" db:"job_name"`
	Orgname    string         `json:"orgname" db:"org_name"`
	RolesList  pq.StringArray `json:"roleList" db:"role_list"`
	Avatar     *Avatar        `json:"avatar" db:"avatar"`
}

type Avatar struct {
	MimeType   string `json:"mimeType"`
	BacketName string `json:"backetName"`
	FileName   string `json:"fileName"`
}

func (s *Avatar) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		err := json.Unmarshal(v, &s)
		if err != nil {
			return err
		}
		return nil
	case string:
		err := json.Unmarshal([]byte(v), &s)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}
