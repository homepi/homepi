package accessory

import (
	"errors"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/gpio/tasks"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) RunAccessory(ctx *gin.Context) {

	var (
		accessory   = new(models.Accessory)
		accessoryId = ctx.Param("accessory_id")
		user        = ctx.MustGet("user").(*models.User)
	)

	if !user.CanRunAccessory() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	result := s.db.Where("id =?", accessoryId).Find(accessory)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(respond.Default.NotFound())
		return
	}

	if result.Error != nil {
		ctx.JSON(respond.Default.SetStatusText("Failed!").
			SetStatusCode(http.StatusInternalServerError).
			RespondWithMessage("Could not get accessory!"))
		return
	}

	if !accessory.IsActive {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusServiceUnavailable).
			SetStatusText("Failed!").
			RespondWithMessage("Accessory is not active!"))
		return
	}

	if _, err := tasks.RunAccessory(s.db, accessory); err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not run accessory. Please contact to your server administrator!"))
		return
	}

	dbLog := &models.Log{Type: models.UserLogType, AccessoryId: accessory.ID, UserId: user.ID}

	if err := s.db.Create(dbLog).Error; err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not run accessory. Please contact to your server administrator!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
		SetStatusText("Success!").
		RespondWithMessage("Accessory ran successfully!"))
	return
}
