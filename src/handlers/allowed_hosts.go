package handlers

import (
	"net/http"
	"strings"

	"github.com/homepi/homepi/src/core"
)

// getHost tries its best to return the request host.
func getHost(r *http.Request) string {
	host := r.Host
	// Slice off any port information.
	if i := strings.Index(host, ":"); i != -1 {
		host = host[:i]
	}
	return host
}

func wrapAllowedHosts(ctx *core.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

			host := getHost(req)

			for _, h := range ctx.Config.AllowedHosts {
				if h == host {
					next.ServeHTTP(w, req)
					return
				}
			}

			w.WriteHeader(http.StatusForbidden)
			return

		})
	}
}
