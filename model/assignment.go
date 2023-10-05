package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shivasaicharanruthala/webapp/errors"
)

type Assignment struct {
	ID                string     `json:"-"`
	AccountID         string     `json:"-"`
	Name              *string    `json:"name"`
	Points            *int       `json:"points"`
	NoOfAttempts      *int       `json:"num_of_attempts"`
	Deadline          *time.Time `json:"deadline"`
	AssignmentCreated time.Time  `json:"-"`
	AssignmentUpdated time.Time  `json:"-"`
}

type AssignmentResponse struct {
	ID                string     `json:"id"`
	AccountID         string     `json:"-"`
	Name              *string    `json:"name"`
	Points            *int       `json:"points"`
	NoOfAttempts      *int       `json:"num_of_attempts"`
	Deadline          *time.Time `json:"deadline"`
	AssignmentCreated time.Time  `json:"assignment_created"`
	AssignmentUpdated time.Time  `json:"assignment_updated"`
}

func (a *Assignment) Validate() error {
	if err := a.ValidateName(); err != nil {
		return err
	}

	if err := a.ValidatePoints(); err != nil {
		return err
	}

	if err := a.ValidateNoOfAttempts(); err != nil {
		return err
	}

	if err := a.ValidateDeadline(); err != nil {
		return err
	}

	return nil
}

func ValidateID(id string) error {
	// Invalid UUID
	var parseErrors errors.InvalidParam

	_, err := uuid.Parse(id)
	if err != nil {
		parseErrors.Param = append(parseErrors.Param, "id")
	}

	if parseErrors.Param != nil {
		return errors.NewInvalidParam(parseErrors)
	}

	return nil
}

func (a *Assignment) ValidateName() error {
	if a.Name == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"name"}})
	}

	return nil
}

func (a *Assignment) ValidatePoints() error {
	// Missing Points
	if a.Points == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"points"}})
	}

	// Invalid range of points
	var parseErrors errors.InvalidParam
	if *a.Points < 1 || *a.Points > 10 {
		parseErrors.Param = append(parseErrors.Param, "points")
	}

	if parseErrors.Param != nil {
		return errors.NewInvalidParam(parseErrors)
	}

	return nil
}

func (a *Assignment) ValidateNoOfAttempts() error {
	// Missing Points
	if a.NoOfAttempts == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"num_of_attempts"}})
	}

	// Invalid range of points
	var parseErrors errors.InvalidParam
	if *a.NoOfAttempts < 1 || *a.NoOfAttempts > 100 {
		parseErrors.Param = append(parseErrors.Param, "num_of_attempts")
	}

	if parseErrors.Param != nil {
		return errors.NewInvalidParam(parseErrors)
	}

	return nil
}

func (a *Assignment) ValidateDeadline() error {
	// Missing Points
	if a.Deadline == nil {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"deadline"}})
	}

	// Invalid range of points
	var parseErrors errors.InvalidParam
	currentTimestamp := time.Now().UTC()
	givenTimestamp := (*a.Deadline).UTC()

	if givenTimestamp.Before(currentTimestamp) || givenTimestamp.Equal(currentTimestamp) {
		parseErrors.Param = append(parseErrors.Param, "deadline")
	}

	if parseErrors.Param != nil {
		return errors.NewInvalidParam(parseErrors)
	}

	return nil
}

func (a *Assignment) SetID() {
	a.ID = uuid.New().String()
}

func (a *Assignment) SetTimestamps(created bool) {
	*a.Deadline = (*a.Deadline).UTC()

	if created {
		a.AssignmentCreated = time.Now().UTC()
	}

	a.AssignmentUpdated = time.Now().UTC()
}

func (a *Assignment) SetAccountID(accountID string) {
	a.AccountID = accountID
}
