package assignment

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shivasaicharanruthala/webapp/errors"
	"github.com/shivasaicharanruthala/webapp/types"
	"io"
	"net/http"

	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/responder"
	"github.com/shivasaicharanruthala/webapp/service"
)

type assignmentService struct {
	ctx               *types.Context
	assignmentService service.Assignment
}

func New(ctx *types.Context, assignmentSvc service.Assignment) *assignmentService {
	return &assignmentService{ctx: ctx, assignmentService: assignmentSvc}
}

func (a *assignmentService) Get(w http.ResponseWriter, r *http.Request, user *model.User) {
	if len(r.URL.Query()) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the request body is not empty
	if len(body) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assignments, err := a.assignmentService.Get(a.ctx, user.ID)
	if err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	responder.SetResponse(assignments, 200, w)
	return
}

func (a *assignmentService) GetById(w http.ResponseWriter, r *http.Request, user *model.User) {
	if len(r.URL.Query()) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the request body is not empty
	if len(body) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	assignmentID := mux.Vars(r)["id"]
	if err = model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	assignment, err := a.assignmentService.GetById(a.ctx, user.ID, assignmentID)
	if err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	responder.SetResponse(assignment, 200, w)
	return
}

func (a *assignmentService) Insert(w http.ResponseWriter, r *http.Request, user *model.User) {
	var assignment model.Assignment

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	assignment.SetAccountID(user.ID)
	assignmentResp, err := a.assignmentService.Insert(a.ctx, &assignment)
	if err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	responder.SetResponse(assignmentResp, 201, w)
	return
}

func (a *assignmentService) Modify(w http.ResponseWriter, r *http.Request, user *model.User) {
	var assignment model.Assignment

	assignmentID := mux.Vars(r)["id"]
	if err := model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	assignment.SetAccountID(user.ID)
	assignment.ID = assignmentID
	_, err = a.assignmentService.Modify(a.ctx, &assignment)
	if err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	w.WriteHeader(204)
	return
}

func (a *assignmentService) Delete(w http.ResponseWriter, r *http.Request, user *model.User) {
	assignmentID := mux.Vars(r)["id"]
	if err := model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	if err := a.assignmentService.Delete(a.ctx, user.ID, assignmentID); err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	w.WriteHeader(204)
	return
}

func (a *assignmentService) PostAssignmentSubmission(w http.ResponseWriter, r *http.Request, user *model.User) {
	var submission model.Submission

	assignmentID := mux.Vars(r)["id"]
	if err := model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	err = json.Unmarshal(body, &submission)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w, r)
		return
	}

	submission.SetAssignmentID(assignmentID)
	submission.SetUserID(user.ID)

	submissionResp, err := a.assignmentService.PostSubmission(a.ctx, &submission)
	if err != nil {
		responder.SetErrorResponse(err, w, r)
		return
	}

	responder.SetResponse(submissionResp, 201, w)
	return
}
