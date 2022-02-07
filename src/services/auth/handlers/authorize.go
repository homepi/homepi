package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"github.com/homepi/homepi/src/validators"
	"github.com/mrjosh/respond.go"
	"gorm.io/gorm"
)

func HandleAuthTokens(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createAuthTokens(ctx).ServeHTTP(w, r)
			return
		case http.MethodPut:
			refreshAuthToken(ctx).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func createAuthTokens(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			user    = new(models.User)
			request = &models.Auth{
				User: r.PostFormValue("user"),
				Pass: r.PostFormValue("pass"),
			}
		)

		if errors := validators.NewValidator(request); len(errors) != 0 {
			respond.NewWithWriter(w).ValidationErrors(errors)
			return
		}

		err := ctx.Database.Where("username =?", request.User).Or("email =?", request.User).Find(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond.NewWithWriter(w).SetStatusCode(http.StatusNotFound).
				SetStatusText("failed").
				RespondWithMessage("User does not exists!")
			return
		}

		if user.ValidatePassword(request.Pass) {

			token, refreshedToken, err := CreateNewTokens(ctx, user)
			if err != nil {
				log.Println(err)
				respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
					SetStatusText("failed").
					RespondWithMessage("Could not create tokens!")
				return
			}

			respond.NewWithWriter(w).Succeed(map[string]interface{}{
				"token":           string(token),
				"refreshed_token": string(refreshedToken),
				"type":            "bearer",
			})
			return
		}

		respond.NewWithWriter(w).SetStatusCode(420).
			SetStatusText("failed").
			RespondWithMessage("Invalid credentials!")

	})
}

func refreshAuthToken(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		tokenString = strings.ReplaceAll(tokenString, "Bearer ", "")

		if tokenString == "" {
			respond.NewWithWriter(w).SetStatusCode(422).
				SetStatusText("failed").
				RespondWithMessage("Token is required!")
			return
		}

		token, refreshedToken, err := refreshToken(ctx, tokenString)
		if err != nil {
			log.Printf("could not refresh token: %v", err)
			respond.NewWithWriter(w).Error(http.StatusBadRequest, 3011)
			return
		}

		respond.NewWithWriter(w).Succeed(map[string]interface{}{
			"token":           string(token),
			"refreshed_token": string(refreshedToken),
			"type":            "bearer",
		})

	})
}
