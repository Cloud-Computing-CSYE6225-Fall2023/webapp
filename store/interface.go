package store

import (
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/types"
)

type Assignment interface {
	Get(ctx *types.Context, userID string) ([]*model.AssignmentResponse, error)
	GetById(ctx *types.Context, assignmentID string) (*model.AssignmentResponse, error)
	IfExists(ctx *types.Context, assignmentID string) (*model.User, error)
	Insert(*types.Context, *model.Assignment) error
	Modify(*types.Context, *model.Assignment) error
	Delete(ctx *types.Context, assignmentID string) error
}

type Account interface {
	Insert(ctx *types.Context, account *model.Account) (*model.Account, error)
	BulkInsert(ctx *types.Context, cols []string, rows [][]string) error
	IsAccountExists(ctx *types.Context, email string) (*model.Account, error)
	FlushData(ctx *types.Context, tableName string) error
}
