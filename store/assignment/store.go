package assignment

import (
	"database/sql"
	cr "errors"
	"github.com/shivasaicharanruthala/webapp/types"

	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/store"
)

type assignmentStore struct {
	DB *sql.DB
}

func New(db *sql.DB) store.Assignment {
	return &assignmentStore{DB: db}
}

func (a *assignmentStore) Get(ctx *types.Context, userID string) ([]*model.AssignmentResponse, error) {
	var assignments = []*model.AssignmentResponse{}

	rows, err := a.DB.Query(GetQuery)
	if err != nil {
		return nil, errors.NewCustomError(err)
	}

	for rows.Next() {
		var assignment model.AssignmentResponse

		err = rows.Scan(&assignment.ID, &assignment.Name, &assignment.Points, &assignment.NoOfAttempts, &assignment.Deadline, &assignment.AssignmentCreated, &assignment.AssignmentUpdated)
		if err != nil {
			return nil, errors.NewCustomError(err)
		}

		assignments = append(assignments, &assignment)
	}

	return assignments, nil
}

func (a *assignmentStore) GetById(ctx *types.Context, assignmentID string) (*model.AssignmentResponse, error) {
	row := a.DB.QueryRow(GetByIDQuery, assignmentID)
	if row.Err() != nil {
		return nil, errors.NewCustomError(row.Err())
	}

	var assignment model.AssignmentResponse
	if err := row.Scan(&assignment.ID, &assignment.AccountID, &assignment.Name, &assignment.Points, &assignment.NoOfAttempts, &assignment.Deadline, &assignment.AssignmentCreated, &assignment.AssignmentUpdated); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.NewEntityNotFound(errors.EntityNotFound{Entity: "assignments", ID: assignmentID})
		}

		return nil, errors.NewCustomError(err)
	}

	return &assignment, nil
}

func (a *assignmentStore) IfExists(ctx *types.Context, assignmentID string) (*model.User, error) {
	row := a.DB.QueryRow(IsAssignmentExistsQuery, assignmentID)
	if row.Err() != nil {
		return nil, errors.NewCustomError(row.Err())
	}

	var user model.User
	if err := row.Scan(&user.ID); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.NewEntityNotFound(errors.EntityNotFound{Entity: "assignments", ID: assignmentID})
		}

		return nil, errors.NewCustomError(err)
	}

	return &user, nil
}

func (a *assignmentStore) Insert(ctx *types.Context, assignment *model.Assignment) error {
	_, err := a.DB.Exec(InsertQuery, assignment.ID, assignment.AccountID, *assignment.Name, *assignment.Points, *assignment.NoOfAttempts, *assignment.Deadline, assignment.AssignmentCreated, assignment.AssignmentUpdated)
	if err != nil {
		return errors.NewCustomError(err)
	}

	return nil
}

func (a *assignmentStore) Modify(ctx *types.Context, assignment *model.Assignment) error {
	res, err := a.DB.Exec(UpdateQuery, *assignment.Name, *assignment.Points, *assignment.NoOfAttempts, *assignment.Deadline, assignment.AssignmentUpdated, assignment.ID)
	if err != nil {
		return errors.NewCustomError(err)
	}

	rowsChanged, err := res.RowsAffected()
	if err != nil {
		return errors.NewCustomError(err)
	}

	if int(rowsChanged) == 0 {
		return errors.NewCustomError(cr.New("Rows not affected."), 400)
	}

	return nil
}

func (a *assignmentStore) Delete(ctx *types.Context, assignmentID string) error {
	res, err := a.DB.Exec(DeleteQuery, assignmentID)
	if err != nil {
		return errors.NewCustomError(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.NewCustomError(err)
	}

	if int(rowsAffected) == 0 {
		return errors.NewCustomError(cr.New("Rows not affected."), 400)
	}

	return nil
}
