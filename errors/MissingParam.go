package errors

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type MissingParam struct {
	Param      []string  `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"400"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewMissingParam(err error) MissingParam {
	return MissingParam{
		Msg:        err.Error(),
		StatusCode: http.StatusBadRequest,
		TimeStamp:  time.Now().UTC(),
	}
}
func (e MissingParam) Error() string {
	if len(e.Param) > 1 {
		return fmt.Sprintf("Parameters " + strings.Join(e.Param, ", ") + " are required for this request")
	} else if len(e.Param) == 1 {
		return fmt.Sprintf("Parameter " + e.Param[0] + " is required for this request")
	} else {
		return "This request is missing parameters"
	}
}
