package handlers

import (
	"net/http"
	"strconv"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
)

func HandleRemoveAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			ctx.WrapAuthAdminHandler(removeAccessory(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func removeAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			accessoryID = r.URL.Query().Get("accessory_id")
			user        = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanRemoveAccessory() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		accID, err := strconv.Atoi(accessoryID)
		if err != nil {
			respond.NewWithWriter(w).SetStatusCode(420).
				SetStatusText("failed").
				RespondWithMessage("AccessoryId is invalid")
			return
		}

		if err = ctx.Database.Where("id =?", accID).Delete(&models.Accessory{}).Error; err != nil {
			respond.NewWithWriter(w).SetStatusText("Failed").
				SetStatusCode(http.StatusBadRequest).
				RespondWithMessage("Could not remove accessory!")
			return
		}

		respond.NewWithWriter(w).SetStatusCode(http.StatusOK).
			SetStatusText("Success!").
			RespondWithMessage("Accessory removed successfully!")

	})
}
