package webhook

import (
	"errors"
	"log"
	"math/big"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) CreateWebhook(ctx *gin.Context) {

	var (
		isPublic = ctx.PostForm("is_public")
		user     = ctx.MustGet("user").(*models.User)
	)

	if !user.CanCreateWebhook() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	accessoryId := new(big.Int)
	accessoryId, ok := accessoryId.SetString(ctx.PostForm("accessory_id"), 10)
	if !ok {
		ctx.JSON(respond.Default.SetStatusText("Failed").
			SetStatusCode(500).
			RespondWithMessage("Internal server error!"))
		return
	}

	webhook := &models.Webhook{
		Name:        ctx.PostForm("name"),
		AccessoryId: uint32(accessoryId.Int64()),
		UserId:      user.ID,
	}

	if isPublic == "1" || isPublic == "true" {
		webhook.IsPublic = true
	}

	result := s.db.Where("id =?", webhook.AccessoryId).
		Find(&webhook.Accessory)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"accessory_id": "Could not find accessory!",
		}))
		return
	}

	if errors := validators.NewValidator(webhook); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	if err := s.db.Create(webhook).Error; err != nil {

		log.Println(err)

		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not create accessory!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(webhook))
	return

}
