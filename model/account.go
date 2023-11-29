package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shivasaicharanruthala/webapp/errors"
)

type Account struct {
	ID             string    `json:"-"`
	FirstName      *string   `json:"first_name"`
	LastName       *string   `json:"last_name"`
	Email          *string   `json:"email"`
	Password       *string   `json:"-"`
	AccountCreated time.Time `json:"-"`
	AccountUpdated time.Time `json:"-"`
}

type User struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}

func (a *Account) Validate() error {
	if err := a.ValidateFirstName(); err != nil {
		return err
	}

	if err := a.ValidateLastName(); err != nil {
		return err
	}

	if err := a.ValidateEmail(); err != nil {
		return err
	}

	if err := a.ValidatePassword(); err != nil {
		return err
	}

	return nil
}

func (a *Account) ValidateID() error {
	// Invalid UUID
	var parseErrors errors.InvalidParam

	_, err := uuid.Parse(a.ID)
	if err != nil {
		parseErrors.Param = append(parseErrors.Param, "id")
	}

	if parseErrors.Param != nil {
		return errors.NewInvalidParam(parseErrors)
	}

	return nil
}

func (a *Account) ValidateFirstName() error {
	if a.FirstName == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"first_name"}})
	}

	return nil
}

func (a *Account) ValidateLastName() error {
	// Missing Points
	if a.LastName == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"last_name"}})
	}

	return nil
}

func (a *Account) ValidateEmail() error {
	// Missing Points
	if a.Email == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"email"}})
	}

	// Invalid range of points
	if err := ValidateEmail(*a.Email); err != nil {
		return err
	}

	return nil
}

func (a *Account) ValidatePassword() error {
	// Missing Points
	if a.Password == nil {
		return errors.MissingParam{Param: []string{"password"}}
	}

	// Error while hashing
	hashedPassword, err := HashPassword(*a.Password)
	if err != nil {
		return err
	}

	*a.Password = hashedPassword

	return nil
}

func (a *Account) SetID() {
	a.ID = uuid.New().String()
}

func (a *Account) SetTimestamps() {
	a.AccountCreated = time.Now().UTC()
	a.AccountUpdated = time.Now().UTC()
}
