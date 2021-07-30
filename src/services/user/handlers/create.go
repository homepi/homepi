package handlers

import (
	"net/http"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	authHandler "github.com/homepi/homepi/src/services/auth/handlers"
	"github.com/homepi/homepi/src/validators"
	"github.com/mrjosh/respond.go"
)

func HandleCreateUser(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			ctx.WrapAuthAdminHandler(createUser(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

// Create a new user through api request
func createUser(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authUser := r.Context().Value(core.ContextUserKey).(*models.User)
		if !authUser.CanCreateUser() {
			respond.NewWithWriter(w).SetStatusCode(403).
				SetStatusText("Failed!").
				RespondWithMessage("Forbidden!")
			return
		}

		var (
			role = models.GetRoleByName(r.PostFormValue("role"))
			user = &models.User{
				Fullname: r.PostFormValue("fullname"),
				Username: r.PostFormValue("username"),
				Email:    r.PostFormValue("email"),
				Password: r.PostFormValue("password"),
				Role:     role,
			}
		)

		if errors := validators.NewValidator(user); len(errors) != 0 {
			respond.NewWithWriter(w).ValidationErrors(errors)
			return
		}

		if r.PostFormValue("password") != r.PostFormValue("password_confirmation") {
			respond.NewWithWriter(w).ValidationErrors(map[string]interface{}{
				"password": []string{
					"Passwords are not match!",
				},
			})
			return
		}

		if err := ctx.Database.Create(user).Error; err != nil {
			respond.NewWithWriter(w).SetStatusCode(420).
				SetStatusText("failed").
				RespondWithMessage("Could not create user.")
			return
		}

		token, refreshedToken, err := authHandler.CreateNewTokens(ctx, user)
		if err != nil {
			respond.NewWithWriter(w).SetStatusCode(http.StatusInternalServerError).
				SetStatusText("Failed!").
				RespondWithMessage("Could not create tokens!")
			return
		}

		respond.NewWithWriter(w).Succeed(map[string]interface{}{
			"token":           string(token),
			"refreshed_token": string(refreshedToken),
			"type":            "bearer",
		})

	})
}
