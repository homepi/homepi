package handlers

import (
	"fmt"
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleListUsers(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthAdminHandler(getAllUsersList(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func HandleListRoles(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthAdminHandler(getAllRolesList(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func getAllRolesList(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var roles = make([]models.Role, 0)
		if err := ctx.Database.Order("created_at desc").Find(&roles).Error; err != nil {
			respond.NewWithWriter(w).
				SetStatusCode(http.StatusInternalServerError).
				SetStatusText("failed").
				RespondWithMessage(fmt.Sprintf("Internal server error: %v", err))
			return
		}
		respond.NewWithWriter(w).Succeed(roles)

	})
}

func getAllUsersList(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var users = make([]models.User, 0)

		// PRELOAD
		err := ctx.Database.
			Preload("Role").
			Order("joined_at desc").
			Find(&users).Error

		if err != nil {
			respond.NewWithWriter(w).
				SetStatusCode(http.StatusInternalServerError).
				SetStatusText("failed").
				RespondWithMessage(fmt.Sprintf("Internal server error: %v", err))
			return
		}

		respond.NewWithWriter(w).Succeed(users)

	})
}
