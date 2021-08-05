package handlers

import (
	"fmt"
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleAccessories(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthAdminHandler(getAllAccessories(ctx)).ServeHTTP(w, r)
			return
		case http.MethodPost:
			ctx.WrapAuthAdminHandler(createAccessory(ctx)).ServeHTTP(w, r)
			return
		case http.MethodDelete:
			ctx.WrapAuthAdminHandler(removeAccessory(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func getAllAccessories(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var accessories = make([]models.Accessory, 0)
		if err := ctx.Database.Order("created_at desc").Find(&accessories).Error; err != nil {
			respond.NewWithWriter(w).
				SetStatusCode(http.StatusInternalServerError).
				SetStatusText("failed").
				RespondWithMessage(fmt.Sprintf("Internal server error: %v", err))
			return
		}
		respond.NewWithWriter(w).Succeed(accessories)

	})
}
