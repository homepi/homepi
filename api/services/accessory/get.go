package accessory

import (
	"errors"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) GetAccessory(ctx *gin.Context) {

	var (
		accessory   = new(models.AccessoryWithUserData)
		accessoryId = ctx.Param("accessory_id")
		user        = ctx.MustGet("user").(*models.User)
	)

	if !user.CanSeeAccessories() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	result := s.db.Table("accessories").
		Where("id =?", accessoryId).
		Find(accessory)

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

	accessory.LoadRelations(s.db)

	ctx.JSON(respond.Default.Succeed(accessory))
	return
}
