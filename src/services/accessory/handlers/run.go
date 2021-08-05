package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/homepi/homepi/pkg/gpio/tasks"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
	"gorm.io/gorm"
)

func HandleRunAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(runAccessory(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func runAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			accessory   = new(models.Accessory)
			accessoryID = chi.URLParam(r, "id")
			user        = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanRunAccessory() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("failed").
				RespondWithMessage("Forbidden!")
			return
		}

		if accessoryID == "" {
			respond.NewWithWriter(w).ValidationErrors(map[string]interface{}{
				"accessory_id": "accessory_id is required!",
			})
			return
		}

		err := ctx.Database.Where("id =?", accessoryID).Find(&accessory).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NewWithWriter(w).NotFound()
			return
		}

		if err != nil {
			respond.NewWithWriter(w).SetStatusText("failed").
				SetStatusCode(http.StatusInternalServerError).
				RespondWithMessage("Could not get accessory!")
			return
		}

		if !accessory.IsActive {
			respond.NewWithWriter(w).SetStatusCode(http.StatusServiceUnavailable).
				SetStatusText("failed").
				RespondWithMessage("Accessory is not active!")
			return
		}

		if _, err := tasks.RunAccessory(accessory); err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("failed").
				RespondWithMessage("Could not run accessory. Please contact server administrator!")
			return
		}

		dbLog := &models.Log{Type: models.UserLogType, AccessoryID: accessory.ID, UserID: user.ID}
		if err := ctx.Database.Create(dbLog).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("failed").
				RespondWithMessage("Could not run accessory. Please contact server administrator!")
			return
		}

		respond.NewWithWriter(w).SetStatusCode(http.StatusOK).
			SetStatusText("Success!").
			RespondWithMessage("Accessory ran successfully!")

	})
}
