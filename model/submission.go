package model

import (
	"github.com/google/uuid"
	"github.com/shivasaicharanruthala/webapp/errors"
	"time"
)

type Submission struct {
	ID            string    `json:"-"`
	UserID        string    `json:"-"`
	AssignmentID  string    `json:"-"`
	SubmissionURL *string   `json:"submission_url"`
	Created       time.Time `json:"-"`
	Updated       time.Time `json:"-"`
}

type SubmissionResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"-"`
	AssignmentID  string    `json:"assignment_id"`
	SubmissionURL string    `json:"submission_url"`
	Created       time.Time `json:"submission_date"`
	Updated       time.Time `json:"assignment_updated"`
}

func (s *Submission) ValidateSubmissionURL() error {
	if s.SubmissionURL == nil || *s.SubmissionURL == "" {
		return errors.NewMissingParam(errors.MissingParam{Param: []string{"submission_url"}})
	}

	return nil
}

func (s *Submission) SetID() {
	s.ID = uuid.New().String()
}

func (s *Submission) SetTimestamps() {
	s.Created = time.Now().UTC()
	s.Updated = time.Now().UTC()
}

func (s *Submission) SetAssignmentID(assignmentID string) {
	s.AssignmentID = assignmentID
}

func (s *Submission) SetUserID(userID string) {
	s.UserID = userID
}

func (s *Submission) ConvertToResponse() *SubmissionResponse {
	return &SubmissionResponse{
		ID:            s.ID,
		UserID:        s.UserID,
		AssignmentID:  s.AssignmentID,
		SubmissionURL: *s.SubmissionURL,
		Created:       s.Created,
		Updated:       s.Updated,
	}
}
