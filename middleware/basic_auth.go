package middleware

import (
	"context"
	"github.com/shivasaicharanruthala/webapp/types"
	"net/http"
	"strings"

	"github.com/shivasaicharanruthala/webapp/model"
	"github.com/shivasaicharanruthala/webapp/responder"
	"github.com/shivasaicharanruthala/webapp/service"
)

type BasicAuth struct {
	ctx        *types.Context
	accountSvc service.Account
	handler    func(w http.ResponseWriter, r *http.Request, user *model.User)
}

func NewBasicAuth(ctx *types.Context, handlerToWrap func(w http.ResponseWriter, r *http.Request, user *model.User), accountSvc service.Account) *BasicAuth {
	return &BasicAuth{ctx: ctx, handler: handlerToWrap, accountSvc: accountSvc}
}

func (ba *BasicAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add custom check to check basic auth
	email, password, err := getEmailAndPassword(r.Header.Get("Authorization"))
	if err != nil {
		responder.SetErrorResponse(err, w, r)

		return
	}

	user, err := ba.accountSvc.IsAccountExists(ba.ctx, email, password)
	if err != nil {
		responder.SetErrorResponse(err, w, r)

		return
	}

	ba.handler(w, r, user)
}

type ctxMsg string

const userID ctxMsg = "userID"
const userEmail ctxMsg = "userEmail"

func BasicAuths(ctx *types.Context, accountSvc service.Account) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if ExemptPath(r) {
				inner.ServeHTTP(w, r)
				return
			}

			// add custom check to check basic auth
			email, password, err := getEmailAndPassword(r.Header.Get("Authorization"))
			if err != nil {
				responder.SetErrorResponse(err, w, r)

				return
			}

			user, err := accountSvc.IsAccountExists(ctx, email, password)
			if err != nil {
				responder.SetErrorResponse(err, w, r)

				return
			}

			ctx := context.WithValue(r.Context(), userID, user.ID)
			ctx = context.WithValue(ctx, userEmail, user.Email)
			*r = *r.Clone(ctx)

			inner.ServeHTTP(w, r)
		})
	}
}

func ExemptPath(r *http.Request) bool {
	return strings.HasSuffix(r.URL.Path, "/healthz")
}
