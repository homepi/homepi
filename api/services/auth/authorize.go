package auth

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) CreateAuthToken(ctx *gin.Context) {

	var (
		user    = new(models.User)
		request = &models.Auth{
			User: ctx.PostForm("user"),
			Pass: ctx.PostForm("pass"),
		}
	)

	if errors := validators.NewValidator(request); len(errors) != 0 {
		ctx.JSON(respond.Default.ValidationErrors(errors))
		return
	}

	result := s.db.Where("username =?", request.User).
		Or("email =?", request.User).Find(user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.JSON(respond.Default.SetStatusCode(http.StatusNotFound).
			SetStatusText("Failed!").
			RespondWithMessage("User does not exists!"))
		return
	}

	if user.ValidatePassword(request.Pass) {

		token, refreshedToken, err := s.CreateNewTokens(user)
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

	ctx.JSON(respond.Default.SetStatusCode(420).
		SetStatusText("Failed!").
		RespondWithMessage("Invalid credentials!"))
	return
}

func (s *Service) RefreshAuthToken(ctx *gin.Context) {

	tokenString := ctx.GetHeader("Authorization")
	tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

	if tokenString == "" {

		ctx.AbortWithStatusJSON(respond.Default.SetStatusCode(422).
			SetStatusText("Failed!").
			RespondWithMessage("Token is required!"))
		return
	}

	token, refreshedToken, err := s.RefreshToken(tokenString)
	if err != nil {
		log.Printf("could not refresh token: %v", err)
		ctx.AbortWithStatusJSON(respond.Default.Error(http.StatusBadRequest, 3011))
		return
	}

	ctx.JSON(respond.Default.Succeed(map[string]interface{}{
		"token":           string(token),
		"refreshed_token": string(refreshedToken),
		"type":            "bearer",
	}))
	return
}
