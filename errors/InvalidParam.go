package errors

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type InvalidParam struct {
	Param      []string  `json:"-"`
	Msg        string    `json:"msg"`
	StatusCode int       `json:"-" default:"400"`
	TimeStamp  time.Time `json:"timestamp"`
}

func NewInvalidParam(err error) InvalidParam {
	return InvalidParam{
		Msg:        err.Error(),
		StatusCode: http.StatusBadRequest,
		TimeStamp:  time.Now().UTC(),
	}
}
func (e InvalidParam) Error() string {
	if len(e.Param) > 1 {
		return fmt.Sprintf("Incorrect value for parameters: " + strings.Join(e.Param, ", "))
	} else if len(e.Param) == 1 {
		return fmt.Sprintf("Incorrect value for parameter: " + e.Param[0])
	} else {
		return "This request has invalid parameters"
	}
}
