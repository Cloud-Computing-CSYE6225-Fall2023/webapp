package service

import (
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/types"
)

type Assignment interface {
	Get(ctx *types.Context, userID string) ([]*model.AssignmentResponse, error)
	GetById(ctx *types.Context, userID, assignmentID string) (*model.AssignmentResponse, error)
	Insert(*types.Context, *model.Assignment) (*model.AssignmentResponse, error)
	Modify(*types.Context, *model.Assignment) (*model.AssignmentResponse, error)
	Delete(ctx *types.Context, userID, assignmentID string) error
}

type Account interface {
	Insert(ctx *types.Context, account *model.Account) (*model.Account, error)
	BulkInsert(ctx *types.Context, filepath string) error
	IsAccountExists(ctx *types.Context, email, pass string) (*model.User, error)
}
