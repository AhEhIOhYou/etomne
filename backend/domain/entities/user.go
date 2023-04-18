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

func (user *User) BeforeUpdate() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	hashPassword, _ := security.Hash(user.Password)
	user.Password = string(hashPassword)
	user.UpdatedAt = time.Now()
}

func (user *User) Prepared() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	hashedBytePass, _ := security.Hash(user.Password)
	user.Password = string(hashedBytePass)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "email email"
			}
		}

	case "login":
		if user.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if user.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgotpassword":
		if user.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if user.Name == "" {
			errorMessages["name_required"] = "name is required"
		}
		if user.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if user.Password != "" && len(user.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}
		if user.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if user.Email != "" {
			if err = checkmail.ValidateFormat(user.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errorMessages
}
