package models

type User struct {
	Id         string `json:"-" db:"id"`
	Login      string `json:"login" binding:"required"`
	Password   string `json:"password" binding:"required""`
	Lastname   string `json:"lastname"`
	Firstname  string `json:"firstname"`
	Middlename string `json:"middlename"`
	Jobname    string `json:"jobname"`
	Orgname    string `json:"orgname"`
	Userid     string `json:"-"`
}
