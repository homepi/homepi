package webhook

import (
	"errors"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/gpio/tasks"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) RunWebhook(ctx *gin.Context) {

	var (
		webhook = new(models.Webhook)
		whId    = ctx.Param("id")
		whHash  = ctx.Param("hash")
	)

	result := s.db.Where("id =?", whId).
		Where("hash =?", whHash).
		Find(webhook)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	s.db.Model(webhook).Association("Accessory")
	s.db.Model(webhook.Accessory).Association("User")

	if !webhook.IsActive {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusServiceUnavailable).
			SetStatusText("Failed!").
			RespondWithMessage("Webhook is not active!"))
		return
	}

	if _, err := tasks.RunAccessory(s.db, &webhook.Accessory); err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not run webhook. Please contact to your server administrator!"))
		return
	}

	dbLOG := &models.Log{
		Type:        models.LogWebHook,
		AccessoryId: webhook.Accessory.ID,
		WebhookId:   webhook.ID,
		UserId:      webhook.Accessory.User.ID,
	}

	if err := s.db.Create(dbLOG).Error; err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not run accessory. Please contact to your server administrator!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
		SetStatusText("Success!").
		RespondWithMessage("Webhook ran successfully!"))
	return
}
