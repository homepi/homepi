package user

import (
	"log"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/pagination"
	"github.com/homepi/homepi/api/db/models"
)

// Get the current token user
func (s *Service) GetAccessories(ctx *gin.Context) {

	var (
		accessories = make([]models.Accessory, 0)
		limit       = pagination.GetLimitParam(ctx)
		page        = pagination.GetPageParam(ctx)
		user        = ctx.MustGet("user").(*models.User)
	)

	if !user.CanSeeAccessories() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	dbResult := s.db.Order("created_at desc").
		Limit(limit).
		Find(&accessories)

	if dbResult.Error != nil {
		log.Println(dbResult.Error)
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not get accessories!"))
		return
	}

	result := pagination.Paginator(accessories, page)

	ctx.JSON(respond.Default.Succeed(result))
	return
}
