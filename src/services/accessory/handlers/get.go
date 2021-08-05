package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
	"gorm.io/gorm"
)

func HandleGetAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getAccessory(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func getAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			accessory = new(models.AccessoryWithUserData)
			id        = chi.URLParam(r, "id")
			user      = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanSeeAccessories() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		err := ctx.Database.Model(&models.Accessory{}).Where("id =?", id).Find(&accessory).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NewWithWriter(w).NotFound()
			return
		}

		if err != nil {
			respond.NewWithWriter(w).SetStatusText("Failed!").
				SetStatusCode(http.StatusInternalServerError).
				RespondWithMessage("Could not get accessory!")
			return
		}

		// LOAD REALATIONS
		accessory.LoadRelations(ctx.Database)

		respond.NewWithWriter(w).Succeed(accessory)
	})
}
