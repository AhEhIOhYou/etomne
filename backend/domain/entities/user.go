package entities

import (
	"github.com/AhEhIOhYou/etomne/backend/infrastructure/security"
	"github.com/badoux/checkmail"
	"html"
	"strings"
	"time"
)

type User struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublicUser struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//TODO create response with this structs
type UserResponse struct {
}

type UserAuth struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	res := make([]interface{}, len(users))
	for index, user := range users {
		res[index] = user.PublicUser()
	}
	return res
}

func (user *User) PublicUser() interface{} {
	return &PublicUser{
		ID:   user.ID,
		Name: user.Name,
	}
}

func (userReq *UserRequest) NewUser() *User {
	return &User{
		Name:     userReq.Name,
		Email:    userReq.Email,
		Password: userReq.Password,
	}
}

func (user *User) BeforeUpdate() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	hashPassword, _ := security.Hash(user.Password)
	user.Password = string(hashPassword)
	user.UpdatedAt = time.Now()
}

func (user *User) Prepare() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	hashedBytePass, _ := security.Hash(user.Password)
	user.Password = string(hashedBytePass)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (userReq *UserRequest) ValidateRequst(action string) string {
	var err error

	switch strings.ToLower(action) {
	case "update":
		if userReq.Email == "" {
			return "email required"
		}
		if userReq.Email != "" {
			if err = checkmail.ValidateFormat(userReq.Email); err != nil {
				return "email format error"
			}
		}
	default:
		if userReq.Password == "" {
			return "password is required"
		}
		if userReq.Email == "" {
			return "email is required"
		}
		if userReq.Email != "" {
			if err = checkmail.ValidateFormat(userReq.Email); err != nil {
				return "please provide a valid email"
			}
		}
	}
	return ""
}
