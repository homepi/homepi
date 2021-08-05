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

func HandleRunWebhook(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			runWebhook(ctx).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func runWebhook(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			id      = chi.URLParam(r, "id")
			hash    = chi.URLParam(r, "hash")
			webhook = new(models.Webhook)
			err     = ctx.Database.Where("id =?", id).Where("hash =?", hash).Find(&webhook).Error
		)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NewWithWriter(w).NotFound()
			return
		}

		ctx.Database.Model(webhook).Association("Accessory")
		ctx.Database.Model(webhook.Accessory).Association("User")

		if !webhook.IsActive {
			respond.NewWithWriter(w).SetStatusCode(http.StatusServiceUnavailable).
				SetStatusText("Failed!").
				RespondWithMessage("Webhook is not active!")
			return
		}

		if _, err := tasks.RunAccessory(webhook.Accessory); err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not run webhook. Please contact server administrator!")
			return
		}

		dbLOG := &models.Log{
			Type:        models.LogWebhook,
			AccessoryID: webhook.Accessory.ID,
			WebhookID:   webhook.ID,
			UserID:      webhook.Accessory.User.ID,
		}

		if err := ctx.Database.Create(dbLOG).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not run accessory. Please contact server administrator!")
			return
		}

		respond.NewWithWriter(w).SetStatusCode(http.StatusOK).
			SetStatusText("Success!").
			RespondWithMessage("Webhook ran successfully!")

	})
}
