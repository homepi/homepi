package handlers

import (
	"net/http"

	"github.com/homepi/homepi/pkg/pagination"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleListLogs(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getLogs(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func getLogs(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			limit = pagination.GetLimitParam(r)
			page  = pagination.GetPageParam(r)
			user  = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		logs, err := models.GetLogs(ctx.Database, user, limit)
		if err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not get logs!")
			return
		}

		result := pagination.Paginator(logs, page)
		respond.NewWithWriter(w).Succeed(result)

	})
}
