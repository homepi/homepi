package handlers

import (
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
)

func HandleUsersMe(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getMe(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func getMe(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.NewWithWriter(w).
			Succeed(r.Context().Value(core.ContextUserKey))
	})
}
