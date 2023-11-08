package middleware

import (
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"

	"github.com/shivasaicharanruthala/webapp/log"
	"github.com/shivasaicharanruthala/webapp/responder"
)

type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *StatusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Logging(logger *log.CustomLogger) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			correlationID := getCorrelationID(r)

			srw := &StatusResponseWriter{ResponseWriter: w}
			defer func(res *StatusResponseWriter, req *http.Request) {
				l := log.Message{
					TraceId:    correlationID,
					Duration:   time.Since(start).Microseconds(),
					Method:     req.Method,
					URI:        req.RequestURI,
					StatusCode: res.status,
				}

				l.ErrorMessage = populateErrorMessage(r, res.status)

				if logger.Logger != nil {
					isServerError := (res.status >= http.StatusBadRequest && res.status <= http.StatusUnavailableForLegalReasons) || (res.status >= http.StatusInternalServerError && res.status <= http.StatusNetworkAuthenticationRequired)

					if isServerError {
						l.Level = "ERROR"
					} else {
						l.Level = "INFO"
					}

					logger.Log(&l)
				}
			}(srw, r)

			inner.ServeHTTP(srw, r)
		})
	}
}

func populateErrorMessage(r *http.Request, statusCode int) string {
	var msg string

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		msg, _ = r.Context().Value(responder.ErrorMessage).(string)
	}

	return msg
}

func getCorrelationID(r *http.Request) string {
	id, _ := uuid.NewUUID()
	s := strings.Split(id.String(), "-")

	correlationID := strings.Join(s, "")

	r.Header.Set("Correlation-ID", correlationID)

	return correlationID
}
