package handlers

import (
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleUserPermissions(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getUserRoles(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func getUserRoles(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(core.ContextUserKey).(*models.User)
		respond.NewWithWriter(w).Succeed(user.Role)
	})
}
