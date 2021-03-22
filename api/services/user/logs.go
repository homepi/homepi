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
func (s *Service) GetLogs(ctx *gin.Context) {

	var (
		limit = pagination.GetLimitParam(ctx)
		page  = pagination.GetPageParam(ctx)
		user  = ctx.MustGet("user").(*models.User)
	)

	logs, err := models.GetLogs(s.db, user, limit)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not get logs!"))
		return
	}

	result := pagination.Paginator(logs, page)

	ctx.JSON(respond.Default.Succeed(result))
	return
}
