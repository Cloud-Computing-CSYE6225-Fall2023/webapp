package middleware

import (
	"fmt"
	"net/http"

	"github.com/shivasaicharanruthala/webapp/types"
)

func APICountMetrics(ctx *types.Context) func(inner http.Handler) http.Handler {
	return func(inner http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := maskPathParams(r.RequestURI)

			ctx.Metrics.Increment(fmt.Sprintf("api_calls.%s", r.Method+"-"+path))

			inner.ServeHTTP(w, r)
		})
	}
}
