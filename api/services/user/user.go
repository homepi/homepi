package user

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/db/models"
)

// Get the current token user
func (s *Service) GetMe(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)
	ctx.JSON(respond.Default.Succeed(user))
	return
}
