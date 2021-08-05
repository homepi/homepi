package handlers

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db/models"
	"gorm.io/gorm"
)

func CreateNewTokens(ctx *core.Context, user *models.User) (token, refreshedToken []byte, err error) {
	//generate the auth token
	token, err = createAuthToken(ctx, user)
	if err != nil {
		return
	}
	// generate the refresh token
	refreshedToken, err = createRefreshToken(ctx, user)
	if err != nil {
		return
	}
	return
}

func createAuthToken(ctx *core.Context, user *models.User) ([]byte, error) {

	authTokenExp := time.Now().Add(time.Minute * time.Duration(ctx.Config.JWT.AccessToken.ExpiresAt)).Unix()

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: authTokenExp,
	})

	// generate the auth token string
	token, err := authJwt.SignedString([]byte(ctx.Config.JWT.AccessToken.Value))
	if err != nil {
		return nil, fmt.Errorf("could not SignedString token: %v", err)
	}
	return []byte(token), nil
}

func createRefreshToken(ctx *core.Context, user *models.User) ([]byte, error) {

	var (
		token = &models.RefreshToken{
			ID:      uuid.New().ID(),
			UserID:  user.ID,
			Valid:   true,
			TokenID: libstr.RandomDigits(30),
		}
	)

	if err := ctx.Database.Create(token).Error; err != nil {
		return nil, err
	}

	refreshTokenExp := time.Now().Add(time.Minute * time.Duration(ctx.Config.JWT.RefreshToken.ExpiresAt)).Unix()

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        token.TokenID,
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: refreshTokenExp,
	})

	// generate the refresh token string
	refreshTokenString, err := refreshJwt.SignedString([]byte(ctx.Config.JWT.RefreshToken.Value))
	if err != nil {
		return nil, fmt.Errorf("could not SignedString RefreshToken: %v", err)
	}
	return []byte(refreshTokenString), nil
}

func checkRefreshToken(ctx *core.Context, tokenID string) (*models.RefreshToken, error) {
	refreshedToken := new(models.RefreshToken)
	err := ctx.Database.
		Preload("User").
		Find(refreshedToken, map[string]interface{}{"token_id": tokenID}).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("could not find refresh token from database, token_id: [%s]", tokenID)
		}
		return nil, err
	}
	// Preload ROLE colum
	if refreshedToken.TokenID == tokenID && refreshedToken.Valid {
		return refreshedToken, nil
	}
	return nil, errors.New("could not find refreshed token")
}

func refreshToken(ctx *core.Context, refreshTokenString string) ([]byte, []byte, error) {

	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, new(jwt.StandardClaims), func(token *jwt.Token) (interface{}, error) {
		return []byte(ctx.Config.JWT.RefreshToken.Value), nil
	})
	if err != nil {
		return nil, nil, err
	}

	if refreshToken == nil {
		return nil, nil, errors.New("error reading jwt claims")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, nil, errors.New("error reading jwt claims")
	}

	var (
		user    = &models.User{}
		findErr = ctx.Database.Where("id =?", refreshTokenClaims.Subject).Find(&user).Error
	)

	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return nil, nil, errors.New("unauthorized")
	}

	dbRefreshedToken, rErr := checkRefreshToken(ctx, refreshTokenClaims.Id)
	if rErr != nil {
		return nil, nil, fmt.Errorf("could not decode refresh token: %v", rErr)
	}

	if !dbRefreshedToken.Valid {
		return nil, nil, errors.New("refresh token is not valid")
	}

	if refreshToken.Valid {
		if err = deleteRefreshToken(ctx, refreshTokenClaims.Id); err != nil {
			return nil, nil, fmt.Errorf("could not delete refresh token: %v", err)
		}
		return CreateNewTokens(ctx, dbRefreshedToken.User)
	}

	return nil, nil, errors.New("unauthorized")
}

func deleteRefreshToken(ctx *core.Context, jti string) (err error) {
	return ctx.Database.Where("token_id =?", jti).Delete(&models.RefreshToken{}).Error
}
