package user

import (
	"log"
	"net/http"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/db/models"
)

// Create a new user through api request
func (s *Service) CreateUser(ctx *gin.Context) {

	authUser := ctx.MustGet("user").(*models.User)
	if !authUser.CanCreateUser() {
		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(403).
			SetStatusText("Failed!").
			RespondWithMessage("Forbidden!"))
		return
	}

	var (
		role = models.GetRoleByName(ctx.PostForm("role"))
		user = &models.User{
			Fullname: ctx.PostForm("fullname"),
			Username: ctx.PostForm("username"),
			Email:    ctx.PostForm("email"),
			Password: ctx.PostForm("password"),
			Role:     *role,
		}
	)

	if errors := validators.NewValidator(user); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	if ctx.PostForm("password") != ctx.PostForm("password_confirmation") {
		ctx.JSON(respond.Default.ValidationErrors(map[string]interface{}{
			"password": []string{
				"Passwords are not match!",
			},
		}))
		return
	}

	if err := s.db.Create(user).Error; err != nil {
		log.Println(err)

		ctx.JSON(respond.Default.SetStatusCode(420).
			SetStatusText("failed").
			RespondWithMessage("Could not create user."))
		return
	}

	token, refreshedToken, err := s.auth.CreateNewTokens(user)
	if err != nil {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusInternalServerError).
			SetStatusText("Failed!").
			RespondWithMessage("Could not create tokens!"))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":           string(token),
		"refreshed_token": string(refreshedToken),
		"type":            "bearer",
	}))
	return
}
