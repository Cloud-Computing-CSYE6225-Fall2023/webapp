package middleware

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shivasaicharanruthala/webapp/types"
	"net/http"
	"strings"
)

func APICountMetrics(ctx *types.Context) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var path string

			if strings.Contains(r.RequestURI, "/v1/assignments") {
				path = "/v1/assignments"

				_, exists := mux.Vars(r)["id"]
				if exists {
					path += "/{id}"
				}
			} else {
				path = r.RequestURI
			}

			ctx.Metrics.Increment(fmt.Sprintf("api_calls.%s", r.Method+"-"+path))

			inner.ServeHTTP(w, r)
		})
	}
}
