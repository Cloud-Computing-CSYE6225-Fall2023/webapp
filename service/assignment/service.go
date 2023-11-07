package assignment

import (
	cr "errors"
	"github.com/shivasaicharanruthala/webapp/types"

	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/service"
	"github.com/shivasaicharanruthala/webapp/store"
)

type dataStore struct {
	assignmentStore store.Assignment
}

func New(adb store.Assignment) service.Assignment {
	return &dataStore{assignmentStore: adb}
}

func (a *dataStore) Get(ctx *types.Context, userID string) ([]*model.AssignmentResponse, error) {
	return a.assignmentStore.Get(ctx, userID)
}

func (a *dataStore) GetById(ctx *types.Context, userID, assignmentID string) (*model.AssignmentResponse, error) {
	assignment, err := a.assignmentStore.GetById(ctx, assignmentID)
	if err != nil {
		return nil, err
	}

	if assignment.AccountID != userID {
		return nil, errors.NewCustomError(cr.New("Logged in user dont have access to fetch this record"), 403)
	}

	return assignment, nil
}

func (a *dataStore) Insert(ctx *types.Context, assignment *model.Assignment) (*model.AssignmentResponse, error) {
	if err := assignment.Validate(); err != nil {
		return nil, err
	}

	assignment.SetID()
	assignment.SetTimestamps(true)

	if err := a.assignmentStore.Insert(ctx, assignment); err != nil {
		return nil, err
	}

	return a.assignmentStore.GetById(ctx, assignment.ID)
}

func (a *dataStore) Modify(ctx *types.Context, assignment *model.Assignment) (*model.AssignmentResponse, error) {
	if err := assignment.Validate(); err != nil {
		return nil, err
	}

	user, err := a.assignmentStore.IfExists(ctx, assignment.ID)
	if err != nil {
		return nil, err
	}

	if user.ID != assignment.AccountID {
		return nil, errors.NewCustomError(cr.New("Logged in user dont have access to fetch this record"), 403)
	}

	assignment.SetTimestamps(false)
	if err = a.assignmentStore.Modify(ctx, assignment); err != nil {
		return nil, err
	}

	return nil, nil
}

func (a *dataStore) Delete(ctx *types.Context, userID, assignmentID string) error {
	user, err := a.assignmentStore.IfExists(ctx, assignmentID)
	if err != nil {
		return err
	}

	if user.ID != userID {
		return errors.NewCustomError(cr.New("Logged in user dont have access to fetch this record"), 403)
	}

	return a.assignmentStore.Delete(ctx, assignmentID)
}
