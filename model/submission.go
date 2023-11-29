package model

import (
	"github.com/google/uuid"
	"github.com/shivasaicharanruthala/webapp/errors"
	"time"
)

type Submission struct {
	ID            string    `json:"-"`
	User          User      `json:"-"`
	AssignmentID  string    `json:"-"`
	SubmissionURL *string   `json:"submission_url"`
	Created       time.Time `json:"-"`
	Updated       time.Time `json:"-"`
}

type SubmissionResponse struct {
	ID            string    `json:"id"`
	User          User      `json:"-"`
	AssignmentID  string    `json:"assignment_id"`
	SubmissionURL string    `json:"submission_url"`
	Created       time.Time `json:"submission_date"`
	Updated       time.Time `json:"submission_updated"`
}

type PublishSubmission struct {
	ID            string    `json:"id"`
	User          User      `json:"user"`
	AssignmentID  string    `json:"assignment_id"`
	SubmissionURL string    `json:"submission_url"`
	Created       time.Time `json:"submission_date"`
	Updated       time.Time `json:"submission_updated"`
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

func (s *Submission) SetUser(user User) {
	s.User = user
}

func (s *Submission) ConvertToResponse() *SubmissionResponse {
	return &SubmissionResponse{
		ID:            s.ID,
		User:          s.User,
		AssignmentID:  s.AssignmentID,
		SubmissionURL: *s.SubmissionURL,
		Created:       s.Created,
		Updated:       s.Updated,
	}
}

func (s *Submission) ConvertToPublishResponse() *PublishSubmission {
	return &PublishSubmission{
		ID:            s.ID,
		User:          s.User,
		AssignmentID:  s.AssignmentID,
		SubmissionURL: *s.SubmissionURL,
		Created:       s.Created,
		Updated:       s.Updated,
	}
}
