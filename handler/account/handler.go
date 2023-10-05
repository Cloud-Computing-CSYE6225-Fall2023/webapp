package account

import (
	"net/http"

	"github.com/shivasaicharanruthala/webapp/service"
)

type accountService struct {
	acntService service.Account
}

func New(acntSvc service.Account) *accountService {
	return &accountService{acntService: acntSvc}
}

func (a *accountService) Insert(w http.ResponseWriter, r *http.Request) {

}
