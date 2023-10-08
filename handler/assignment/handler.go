package assignment

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/shivasaicharanruthala/webapp/errors"
	"io"
	"net/http"

	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/responder"
	"github.com/shivasaicharanruthala/webapp/service"
)

type assignmentService struct {
	assignmentService service.Assignment
}

func New(assignmentSvc service.Assignment) *assignmentService {
	return &assignmentService{assignmentService: assignmentSvc}
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

	assignments, err := a.assignmentService.Get(user.ID)
	if err != nil {
		responder.SetErrorResponse(err, w)
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
		responder.SetErrorResponse(err, w)
		return
	}

	assignment, err := a.assignmentService.GetById(user.ID, assignmentID)
	if err != nil {
		responder.SetErrorResponse(err, w)
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
		responder.SetErrorResponse(err, w)
		return
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w)
		return
	}

	assignment.SetAccountID(user.ID)
	assignmentResp, err := a.assignmentService.Insert(&assignment)
	if err != nil {
		responder.SetErrorResponse(err, w)
		return
	}

	responder.SetResponse(assignmentResp, 201, w)
	return
}

func (a *assignmentService) Modify(w http.ResponseWriter, r *http.Request, user *model.User) {
	var assignment model.Assignment

	assignmentID := mux.Vars(r)["id"]
	if err := model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w)
		return
	}

	err = json.Unmarshal(body, &assignment)
	if err != nil {
		err = errors.NewCustomError(err, 400)
		responder.SetErrorResponse(err, w)
		return
	}

	assignment.SetAccountID(user.ID)
	assignment.ID = assignmentID
	_, err = a.assignmentService.Modify(&assignment)
	if err != nil {
		responder.SetErrorResponse(err, w)
		return
	}

	w.WriteHeader(204)
	return
}

func (a *assignmentService) Delete(w http.ResponseWriter, r *http.Request, user *model.User) {
	assignmentID := mux.Vars(r)["id"]
	if err := model.ValidateID(assignmentID); err != nil {
		responder.SetErrorResponse(err, w)
		return
	}

	if err := a.assignmentService.Delete(user.ID, assignmentID); err != nil {
		responder.SetErrorResponse(err, w)
		return
	}

	w.WriteHeader(204)
	return
}
