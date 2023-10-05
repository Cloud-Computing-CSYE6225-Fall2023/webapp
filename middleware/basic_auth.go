package middleware

import (
	"net/http"

	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/responder"
	"github.com/shivasaicharanruthala/webapp/service"
)

type BasicAuth struct {
	accountSvc service.Account
	handler    func(w http.ResponseWriter, r *http.Request, user *model.User)
}

func NewBasicAuth(handlerToWrap func(w http.ResponseWriter, r *http.Request, user *model.User), accountSvc service.Account) *BasicAuth {
	return &BasicAuth{handler: handlerToWrap, accountSvc: accountSvc}
}

func (ba *BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add custom check to check basic auth
	email, password, err := getEmailAndPassword(r.Header.Get("Authorization"))
	if err != nil {
		responder.SetErrorResponse(err, w)

		return
	}

	user, err := ba.accountSvc.IsAccountExists(email, password)
	if err != nil {
		responder.SetErrorResponse(err, w)

		return
	}

	ba.handler(w, r, user)
}
