package handlers

import (
	"net/http"
	"strings"

	"github.com/homepi/homepi/pkg/libstr"
	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
)

var allowedMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
}

func wrapCORSHandler(ctx *core.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			corsConf := ctx.CORSConfig()

			// If CORS is not enabled or if no Origin header is present (i.e. the request
			// is from the Vault CLI. A browser will always send an Origin header), then
			// just return a 204.
			if !corsConf.Enabled {
				next.ServeHTTP(w, req)
				return
			}

			origin := req.Header.Get("Origin")
			requestMethod := req.Header.Get("Access-Control-Request-Method")

			if origin == "" {
				next.ServeHTTP(w, req)
				return
			}

			// Return a 403 if the origin is not allowed to make cross-origin requests.
			if !corsConf.IsValidOrigin(origin) {

				respond.NewWithWriter(w).
					SetStatusCode(http.StatusForbidden).
					SetStatusText("failed").
					RespondWithMessage("origin not allowed")
				return
			}

			if req.Method == http.MethodOptions && !libstr.StrListContains(allowedMethods, requestMethod) {

				respond.NewWithWriter(w).MethodNotAllowed()
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")

			// apply headers for preflight requests
			if req.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConf.AllowedHeaders, ","))
				w.Header().Set("Access-Control-Max-Age", "300")
				return
			}

			next.ServeHTTP(w, req)

		})
	}
}
