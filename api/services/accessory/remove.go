package accessory

import (
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/db/models"
)

func (s *Service) RemoveAccessory(ctx *gin.Context) {

	var (
		accessory   = new(models.Accessory)
		accessoryId = ctx.Param("accessory_id")
		user        = ctx.MustGet("user").(*models.User)
	)

	if !user.CanRemoveAccessory() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	result := s.db.Delete(accessory, "id =?", accessoryId)

	if result.Error != nil {
		ctx.JSON(respond.Default.SetStatusText("Failed").
			SetStatusCode(http.StatusBadRequest).
			RespondWithMessage("Could not remove accessory!"))
		return
	}

	ctx.JSON(respond.Default.SetStatusCode(http.StatusOK).
		SetStatusText("Success!").
		RespondWithMessage("Accessory removed successfully!"))
	return
}
