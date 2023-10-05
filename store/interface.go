package store

import "github.com/shivasaicharanruthala/webapp/model"

type Assignment interface {
	Get(userID string) ([]*model.AssignmentResponse, error)
	GetById(assignmentID string) (*model.AssignmentResponse, error)
	IfExists(assignmentID string) (*model.User, error)
	Insert(*model.Assignment) error
	Modify(*model.Assignment) error
	Delete(assignmentID string) error
}

type Account interface {
	Insert(account *model.Account) (*model.Account, error)
	BulkInsert(cols []string, rows [][]string) error
	IsAccountExists(email string) (*model.Account, error)
	FlushData(tableName string) error
}
