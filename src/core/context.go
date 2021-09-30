package core

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/homepi/homepi/src/db/models"
	"github.com/mrjosh/respond.go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ContextUserKeyType string

const (
	ContextUserKey ContextUserKeyType = "user"
)

type Context struct {
	Database *gorm.DB
	Config   *ConfMap
}

func (ctx *Context) authenticate(w http.ResponseWriter, r *http.Request) (*models.User, error) {

	accessToken := strings.TrimSpace(r.Header.Get("Authorization"))
	if accessToken == "" {
		errmsg := "Token is required!"
		respond.NewWithWriter(w).SetStatusCode(422).
			SetStatusText("failed").
			RespondWithMessage(errmsg)
		return nil, errors.New(errmsg)
	}

	tokenSlice := strings.Split(accessToken, " ")
	if len(tokenSlice) < 2 {
		errmsg := "Token is required"
		respond.NewWithWriter(w).SetStatusCode(422).
			SetStatusText("failed").
			RespondWithMessage(errmsg)
		return nil, errors.New(errmsg)
	}

	var (
		user      *models.User
		err       error
		ok        bool
		tokenType = strings.TrimSpace(tokenSlice[0])
		token     = strings.TrimSpace(tokenSlice[1])
	)

	switch tokenType {
	case "Bearer", "bearer":
		user, ok, err = ctx.DecodeAuthToken(token)
	case "ApiToken", "api_token", "apitoken":
		user, ok, err = ctx.CheckAPIToken(token)
	default:
		respond.NewWithWriter(w).
			SetStatusCode(http.StatusBadRequest).
			SetStatusText("failed").
			RespondWithMessage("Token type is invalid")
		return nil, errors.New("Token type is invalid")
	}

	if err != nil {
		respond.NewWithWriter(w).Error(http.StatusUnauthorized, 3011)
		return nil, err
	}

	if !ok {
		respond.NewWithWriter(w).Error(http.StatusUnauthorized, 3011)
		return nil, errors.New("Token is invalid")
	}

	if !user.IsActive {
		respond.NewWithWriter(w).Error(http.StatusBadRequest, 3012)
		return nil, errors.New("User is inactive")
	}

	return user, nil
}

func (ctx *Context) WrapAuthAdminHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := ctx.authenticate(w, r)
		if err != nil {
			logrus.WithError(err).Log(logrus.ErrorLevel)
			return
		}
		rctx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(rctx))
	})
}

// Authenticate given user's token
func (ctx *Context) WrapAuthUserHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := ctx.authenticate(w, r)
		if err != nil {
			logrus.WithError(err).Log(logrus.ErrorLevel)
			//respond.NewWithWriter(w).
			//SetStatusText("failed").
			//SetStatusCode(http.StatusUnauthorized).
			//RespondWithMessage("Unauthorized!")
			return
		}
		rctx := context.WithValue(r.Context(), ContextUserKey, user)
		next.ServeHTTP(w, r.WithContext(rctx))
	})
}

func (ctx *Context) CheckAPIToken(token string) (*models.User, bool, error) {
	var (
		apiToken = models.APIToken{}
		err      = ctx.Database.
				Where("api_tokens.token =?", token).
				Preload("User").
				Preload("Role").
				First(&apiToken).
				Error
	)
	if err != nil {
		return nil, false, fmt.Errorf("error on getting api token: %v", err)
	}
	//apiToken.Role = role
	return apiToken.User, true, nil
}

func (ctx *Context) DecodeAuthToken(token string) (*models.User, bool, error) {

	// now, check that it matches what's in the auth token claims
	authToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ctx.Config.JWT.AccessToken.Value), nil
	})
	if err != nil {
		return nil, false, fmt.Errorf("could not parse jwt.Token with Claims: %v", err)
	}
	if authToken == nil || authToken.Claims == nil {
		return nil, false, errors.New("error reading jwt claims")
	}
	authTokenClaims, ok := authToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, false, errors.New("error reading jwt claims")
	}
	if !authToken.Valid {
		return nil, false, errors.New("auth token is not valid")
	}
	var (
		user    = &models.User{}
		findErr = ctx.Database.
			Where("id =?", authTokenClaims.Subject).
			Preload("Role").
			First(&user).
			Error
	)
	if findErr != nil {
		return nil, false, findErr
	}
	return user, true, nil
}
