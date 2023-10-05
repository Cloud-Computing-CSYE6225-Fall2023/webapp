package service

import "github.com/shivasaicharanruthala/webapp/model"

type Assignment interface {
	Get(userID string) ([]*model.AssignmentResponse, error)
	GetById(userID, assignmentID string) (*model.AssignmentResponse, error)
	Insert(*model.Assignment) (*model.AssignmentResponse, error)
	Modify(*model.Assignment) (*model.AssignmentResponse, error)
	Delete(userID, assignmentID string) error
}

type Account interface {
	Insert(account *model.Account) (*model.Account, error)
	BulkInsert(filepath string) error
	IsAccountExists(email, pass string) (*model.User, error)
}
