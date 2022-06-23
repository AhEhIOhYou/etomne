package entities

type User struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Pass  string `json:"hashedPass"`
	Name  string `json:"name"`
}
