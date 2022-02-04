package handlers

import (
	"net/http"
	"strconv"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/homepi/homepi/src/validators"
	"github.com/mrjosh/respond.go"
	"github.com/stianeikeland/go-rpio"
)

// create accessory handler
func HandleCreateAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ctx.WrapAuthAdminHandler(createAccessory(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func createAccessory(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			isPublic = r.PostFormValue("is_public")
			pinID, _ = strconv.Atoi(r.PostFormValue("pin_id"))

			user = r.Context().Value(core.ContextUserKey).(*models.User)

			taskInt, _ = strconv.Atoi(r.PostFormValue("task"))
			task       = models.Task(taskInt)

			accessory = &models.Accessory{
				Name:        r.PostFormValue("name"),
				Task:        task,
				Description: r.PostFormValue("description"),
				IsPublic:    false,
				PinID:       uint64(pinID),
				User:        user,
				UserID:      user.ID,
				State:       rpio.Low,
				IsActive:    true,
			}
		)

		if !user.CanCreateAccessory() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		if isPublic == "1" || isPublic == "true" {
			accessory.IsPublic = true
		}

		if errors := validators.NewValidator(accessory); len(errors) != 0 {
			respond.NewWithWriter(w).ValidationErrors(errors)
			return
		}

		switch accessory.Task {
		case models.TaskDoor:
			accessory.Icon = "doorbell"
		case models.TaskLamp:
			accessory.Icon = "lamp"
		case models.TaskToggle:
			accessory.Icon = "switch"
		}

		if err := ctx.Database.Create(accessory).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not create accessory!")
			return
		}

		respond.NewWithWriter(w).Succeed(accessory)

	})
}
