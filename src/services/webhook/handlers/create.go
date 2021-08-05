package handlers

import (
	"errors"
	"math/big"
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/homepi/homepi/src/validators"
	"github.com/mrjosh/respond.go"
	"gorm.io/gorm"
)

func HandleWebhooks(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthUserHandler(getWebhooks(ctx)).ServeHTTP(w, r)
			return
		case http.MethodPost:
			ctx.WrapAuthUserHandler(createWebhook(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Get the current token user
func createWebhook(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			isPublic = r.PostFormValue("is_public")
			user     = r.Context().Value(core.ContextUserKey).(*models.User)
		)

		if !user.CanCreateWebhook() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		accessoryID := new(big.Int)
		accessoryID, ok := accessoryID.SetString(r.PostFormValue("accessory_id"), 10)
		if !ok {
			respond.NewWithWriter(w).SetStatusText("Failed").
				SetStatusCode(500).
				RespondWithMessage("Internal server error!")
			return
		}

		webhook := &models.Webhook{
			Name:        r.PostFormValue("name"),
			AccessoryID: uint32(accessoryID.Int64()),
			UserID:      user.ID,
		}

		if isPublic == "1" || isPublic == "true" {
			webhook.IsPublic = true
		}

		err := ctx.Database.Where("id =?", webhook.AccessoryID).Find(&webhook.Accessory).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NewWithWriter(w).ValidationErrors(map[string]interface{}{
				"accessory_id": "Could not find accessory!",
			})
			return
		}

		if errors := validators.NewValidator(webhook); len(errors) != 0 {
			respond.NewWithWriter(w).ValidationErrors(errors)
			return
		}

		if err := ctx.Database.Create(webhook).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not create accessory!")
			return
		}

		respond.NewWithWriter(w).Succeed(webhook)

	})
}
