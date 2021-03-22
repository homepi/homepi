package webhook

import (
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/components/pagination"
	"github.com/homepi/homepi/api/db/models"
)

func (s *Service) GetWebhooks(ctx *gin.Context) {

	var (
		limit = pagination.GetLimitParam(ctx)
		page  = pagination.GetPageParam(ctx)
		user  = ctx.MustGet("user").(*models.User)
	)

	if !user.CanSeeWebhook() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	webhooks, err := models.GetWebhooks(s.db, user, limit)
	if err != nil {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not get logs!"))
		return
	}

	result := pagination.Paginator(webhooks, page)

	ctx.JSON(respond.Default.Succeed(result))
	return
}
