package user

import (
	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/db/models"
)

func (s *Service) SearchUser(ctx *gin.Context) {

	var (
		users   = make([]models.User, 0)
		user    = ctx.MustGet("user").(*models.User)
		request = &models.SearchUser{
			Query: ctx.PostForm("query"),
		}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	if !user.CanSeeUsers() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	s.db.Where("username LIKE ?", "%"+request.Query+"%").
		Or("fullname LIKE ?", "%"+request.Query+"%").
		Or("email LIKE ?", "%"+request.Query+"%").
		Not("id = ?", user.ID).
		Find(&users)

	ctx.JSON(respond.Default.Succeed(users))
	return
}
