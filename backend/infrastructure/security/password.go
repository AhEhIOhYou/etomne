package security

import (
	"fmt"

	"github.com/AhEhIOhYou/etomne/backend/constants"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf(constants.PasswordHashError, err)
	}

	return string(hashedPass), nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf(constants.PasswordVerifyError, err)
	}

	return nil
}
