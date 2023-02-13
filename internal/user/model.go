package user

type Book struct {
	ID         string `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	JobName    string `json:"jobName"`
	OrgName    string `json:"orgName"`
}
