package middleware

import (
	"fmt"
	"github.com/shivasaicharanruthala/webapp/types"
	"net/http"
)

func APICountMetrics(ctx *types.Context) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx.Metrics.Increment(fmt.Sprintf("api_calls.%s", r.Method+"-"+r.RequestURI))

			inner.ServeHTTP(w, r)
		})
	}
}
