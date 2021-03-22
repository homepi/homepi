package user

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/db/models"
)

func (s *Service) GetAllUsers(ctx *gin.Context) {

	user := ctx.MustGet("user").(*models.User)
	if !user.CanSeeUsers() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	users := make([]models.User, 0)
	s.db.Order("joined_at desc").Find(&users)

	for index := range users {
		users[index].LoadRelations(s.db)
	}

	ctx.JSON(respond.Default.Succeed(users))
	return
}
