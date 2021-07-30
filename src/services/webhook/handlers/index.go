package handlers

import (
	"net/http"

	"github.com/homepi/homepi/pkg/pagination"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func getWebhooks(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			limit = pagination.GetLimitParam(r)
			page  = pagination.GetPageParam(r)
			user  = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanSeeWebhook() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		webhooks, err := models.GetWebhooks(ctx.Database, user, limit)
		if err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not get logs!")
			return
		}

		result := pagination.Paginator(webhooks, page)
		respond.NewWithWriter(w).Succeed(result)

	})
}
