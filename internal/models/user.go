package models

import (
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
)

type User struct {
	Id        string `json:"-" db:"id"`
	Login     string `json:"login" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Job       string `json:"job"`
	Org       string `json:"org"`
	Userid    string `json:"-"`
	Role      string `json:"role" db:"name"`
}

type UserList struct {
	Id          string         `json:"id" db:"id"`
	Login       string         `json:"login" db:"login"`
	Lastname    string         `json:"lastname" db:"last_name"`
	Firstname   string         `json:"firstname" db:"first_name"`
	Surname     string         `json:"surname" db:"sur_name"`
	Job         string         `json:"job" db:"job_name"`
	Org         string         `json:"org" db:"org_name"`
	RolesList   pq.StringArray `json:"roleList" db:"role_list"`
	Permissions pq.StringArray `json:"permissions" db:"permissions"`
	Avatar      *Avatar        `json:"avatar" db:"avatar"`
}

type Avatar struct {
	MimeType   string `json:"mimeType" db:"mime_type"`
	BucketName string `json:"bucketName" db:"bucket_name"`
	FileName   string `json:"fileName" db:"file_name"`
}

type Roles struct {
	Value string `json:"value" db:"id"`
	Label string `json:"label" db:"name"`
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

type SingleAvatars struct {
	MimeType   string `json:"mimeType" db:"mime_type"`
	BucketName string `json:"bucketName" db:"bucket_name"`
	FileName   string `json:"fileName" db:"file_name"`
	ImgId      string `json:"_" db:"id"`
}

type UserShortList struct {
	Id        string  `json:"id" db:"id"`
	Lastname  string  `json:"lastname" db:"last_name"`
	Firstname string  `json:"firstname" db:"first_name"`
	Surname   string  `json:"surname" db:"sur_name"`
	Avatar    *Avatar `json:"avatar" db:"avatar"`
}
