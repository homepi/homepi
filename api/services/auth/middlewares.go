package auth

import (
	"net/http"
	"strings"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/db/models"
)

// Authenticate given user's token
func (s *Service) Authentication(ctx *gin.Context) {

	var (
		user  = new(models.User)
		token = ctx.GetHeader("Authorization")
	)

	token = strings.ReplaceAll(token, "Bearer ", "")

	if token == "" {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	user, err := s.DecodeAuthToken(token)
	if err != nil {
		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusUnauthorized, 3011))
		return
	}

	if !user.IsActive {
		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusBadRequest, 3012))
		return
	}

	// Load role relation for user
	user.LoadRelations(s.db)

	// set user in gin.Context
	ctx.Set("user", user)
	ctx.Next()
}
