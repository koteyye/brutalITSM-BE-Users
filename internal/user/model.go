package user

type Person struct {
	ID         string `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	JobName    string `json:"jobName"`
	OrgName    string `json:"orgName"`
	UserId     string `json:"userId"`
}

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
