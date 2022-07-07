package entities

import (
	"html"
	"strings"
)

type User struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"hashedPass"`
	Name     string `json:"name"`
}

type PublicUser struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (u *User) BeforeSave() error {
	hashPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

type Users []User

func (users Users) PublicUser() []interface{} {
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

func (user *User) Prepared() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
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
