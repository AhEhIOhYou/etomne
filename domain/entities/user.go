package entities

import (
	"etomne/infrastructure/security"
	"github.com/badoux/checkmail"
	"html"
	"strings"
	"time"
)

type User struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"size:100;not null;" json:"name"`
	Email     string     `gorm:"size:100;not null;unique" json:"email"`
	Password  string     `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type PublicUser struct {
	ID   uint64 `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"size:100;not null;" json:"name"`
}

func (user *User) BeforeSave() error {
	hashPassword, err := security.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashPassword)
	return nil
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

func (user *User) Prepared() {
	user.Name = html.EscapeString(strings.TrimSpace(user.Name))
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
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
