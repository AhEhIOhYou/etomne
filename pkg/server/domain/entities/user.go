package entities

import (
	"html"
	"strings"
	"time"

	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
	"github.com/AhEhIOhYou/etomne/pkg/server/infrastructure/security"
	"github.com/badoux/checkmail"
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
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuth struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type UserResponse struct {
	PublicUser `json:"public_data"`
	UserAuth   `json:"tokens"`
}

type Users []User

func (users Users) PublicUsers() []interface{} {
	res := make([]interface{}, len(users))
	for index, user := range users {
		res[index] = user.PublicUser()
	}
	return res
}

func (user *User) PublicUser() *PublicUser {
	return &PublicUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
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

func (user *LoginRequest) Prepare() {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
}

func (user UserRequest) Validate() string {
	if user.Password == "" {
		return constants.PasswordCantBeEmpty
	}
	if user.Email == "" {
		return constants.EmailCantBeEmpty
	} else if err := checkmail.ValidateFormat(user.Email); err != nil {
		return constants.EmailWrongFormat
	}
	return ""
}

func (user User) Validate() string {
	if user.Password == "" {
		return constants.PasswordCantBeEmpty
	}
	if user.Email == "" {
		return constants.EmailCantBeEmpty
	} else if err := checkmail.ValidateFormat(user.Email); err != nil {
		return constants.EmailWrongFormat
	}
	return ""
}
