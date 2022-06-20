package entities

type Model struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CreateDate  string `json:"createDate"`
	Description string `json:"description"`
}
