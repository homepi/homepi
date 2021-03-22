package accessory

import (
	"log"
	"net/http"
	"strconv"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/db/models"
	"github.com/stianeikeland/go-rpio"
)

func (s *Service) CreateAccessory(ctx *gin.Context) {

	var (
		isPublic = ctx.PostForm("is_public")
		pinId, _ = strconv.Atoi(ctx.PostForm("pin_id"))

		user = ctx.MustGet("user").(*models.User)

		taskInt, _ = strconv.Atoi(ctx.PostForm("task"))
		task       = models.Task(taskInt)

		accessory = &models.Accessory{
			Name:        ctx.PostForm("name"),
			Task:        task,
			Description: ctx.PostForm("description"),
			IsPublic:    false,
			PinId:       uint64(pinId),
			User:        *user,
			UserId:      user.ID,
			State:       rpio.Low,
		}
	)

	if !user.CanCreateAccessory() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	if isPublic == "1" || isPublic == "true" {
		accessory.IsPublic = true
	}

	if errors := validators.NewValidator(accessory); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	switch accessory.Task {
	case models.TaskDoor:
		accessory.Icon = "doorbell"
	case models.TaskLamp:
		accessory.Icon = "lamp"
	case models.TaskToggle:
		accessory.Icon = "switch"
	case models.TaskFan:
		accessory.Icon = "fan"
	}

	if err := s.db.Create(accessory).Error; err != nil {

		log.Println(err)

		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not create accessory!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(accessory))
	return
}
