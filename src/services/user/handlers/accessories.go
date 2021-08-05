package handlers

import (
	"net/http"

	"github.com/homepi/homepi/pkg/pagination"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleListAccessories(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getUserAccessories(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func getUserAccessories(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			accessories = make([]models.Accessory, 0)
			limit       = pagination.GetLimitParam(r)
			page        = pagination.GetPageParam(r)
			user        = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanSeeAccessories() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		if err := ctx.Database.Order("created_at desc").Limit(limit).Find(&accessories).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not get accessories!")
			return
		}

		result := pagination.Paginator(accessories, page)
		respond.NewWithWriter(w).Succeed(result)

	})
}
