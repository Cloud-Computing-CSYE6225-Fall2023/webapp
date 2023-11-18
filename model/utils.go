package model

import (
	"github.com/google/uuid"
	"regexp"

	"github.com/shivasaicharanruthala/webapp/errors"
	"golang.org/x/crypto/bcrypt"
)

func ValidateEmail(email string) error {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~[:^ascii:]-]+@(?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9]\\.)+[a-zA-Z]{2,}$")

	if !emailRegex.MatchString(email) {
		return errors.NewInvalidParam(errors.InvalidParam{Param: []string{"email"}})
	}

	return nil
}

// VerifyPassword verifies a password against a hashed password and returns an error if verification fails
func VerifyPassword(hashedPassword, enteredPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword))
	return err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.NewInvalidParam(errors.InvalidParam{Param: []string{"password"}})
	}

	return string(hashedPassword), nil
}

func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
