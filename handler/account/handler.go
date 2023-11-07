package account

import (
	"github.com/shivasaicharanruthala/webapp/types"
	"net/http"

	"github.com/shivasaicharanruthala/webapp/service"
)

type accountService struct {
	ctx         *types.Context
	acntService service.Account
}

func New(ctx *types.Context, acntSvc service.Account) *accountService {
	return &accountService{ctx: ctx, acntService: acntSvc}
}

func (a *accountService) Insert(w http.ResponseWriter, r *http.Request) {

}
