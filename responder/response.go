package responder

import (
	"encoding/json"
	"net/http"

	"github.com/shivasaicharanruthala/webapp/errors"
)

func SetErrorResponse(err error, w http.ResponseWriter) {
	switch val := err.(type) {
	case errors.InvalidParam:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.MissingParam:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.EntityNotFound:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	case errors.CustomError:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(val.StatusCode)
		_, _ = w.Write(errJson)
	default:
		errJson, _ := json.Marshal(val)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		_, _ = w.Write(errJson)
	}
}

func SetResponse(resp interface{}, statusCode int, w http.ResponseWriter) {
	respJson, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(respJson)
}
