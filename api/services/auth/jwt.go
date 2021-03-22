package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/homepi/homepi/api/db/models"
	"gorm.io/gorm"
)

func (s *Service) CreateNewTokens(user *models.User) (token, refreshedToken []byte, err error) {
	//generate the auth token
	token, err = s.createAuthToken(user)
	if err != nil {
		return
	}
	// generate the refresh token
	refreshedToken, err = s.createRefreshToken(user)
	if err != nil {
		return
	}
	return
}

func (s *Service) createAuthToken(user *models.User) ([]byte, error) {

	authTokenExp := time.Now().Add(time.Minute * s.jwt.AccessToken.ExpireTime).Unix()

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: authTokenExp,
	})

	// generate the auth token string
	token, err := authJwt.SignedString(s.jwt.AccessToken.Secret)
	if err != nil {
		return nil, fmt.Errorf("could not SignedString token: %v", err)
	}
	return []byte(token), nil
}

func (s *Service) createRefreshToken(user *models.User) ([]byte, error) {

	var (
		token = &models.RefreshedToken{
			UserId: user.ID,
			Valid:  true,
		}
	)

	if err := s.db.Create(token).Error; err != nil {
		return nil, err
	}

	refreshTokenExp := time.Now().Add(time.Minute * s.jwt.RefreshToken.ExpireTime).Unix()

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        token.TokenId,
		Subject:   fmt.Sprint(user.ID),
		ExpiresAt: refreshTokenExp,
	})

	// generate the refresh token string
	refreshTokenString, err := refreshJwt.SignedString(s.jwt.RefreshToken.Secret)
	if err != nil {
		return nil, fmt.Errorf("could not SignedString RefreshToken: %v", err)
	}
	return []byte(refreshTokenString), nil
}

func (s *Service) checkRefreshToken(tokenID string) (*models.RefreshedToken, error) {
	refreshedToken := new(models.RefreshedToken)
	result := s.db.First(refreshedToken, map[string]interface{}{"token_id": tokenID})
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("could not find refresh token from database, token_id: [%s]", tokenID)
		}
		return nil, err
	}
	s.db.Model(refreshedToken).Association("User")
	if refreshedToken.TokenId == tokenID && refreshedToken.Valid {
		return refreshedToken, nil
	}
	return nil, errors.New("could not find refreshed token")
}

func (s *Service) RefreshToken(refreshTokenString string) ([]byte, []byte, error) {

	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, new(jwt.StandardClaims), func(token *jwt.Token) (interface{}, error) {
		return s.jwt.RefreshToken.Secret, nil
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
		user = new(models.User)
		res  = s.db.First(&user, map[string]interface{}{"hash": refreshTokenClaims.Subject})
	)

	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, nil, errors.New("unauthorized")
	}

	dbRefreshedToken, rErr := s.checkRefreshToken(refreshTokenClaims.Id)
	if rErr != nil {
		return nil, nil, fmt.Errorf("could not decode refresh token: %v", rErr)
	}

	if !dbRefreshedToken.Valid {
		return nil, nil, errors.New("refresh token is not valid")
	}

	if refreshToken.Valid {
		if err = s.deleteRefreshToken(refreshTokenClaims.Id); err != nil {
			return nil, nil, fmt.Errorf("could not delete refresh token: %v", err)
		}
		return s.CreateNewTokens(&dbRefreshedToken.User)
	}

	return nil, nil, errors.New("unauthorized")
}

func (s *Service) deleteRefreshToken(jti string) (err error) {
	result := s.db.Delete(&models.RefreshedToken{TokenId: jti})
	if result.Error != nil {
		return errors.New("could not delete refresh token from db")
	}
	return nil
}

func (s *Service) DecodeAuthToken(token string) (*models.User, error) {
	// now, check that it matches what's in the auth token claims
	authToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwt.AccessToken.Secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse jwt.Token with Claims: %v", err)
	}
	if authToken == nil || authToken.Claims == nil {
		return nil, errors.New("error reading jwt claims")
	}
	authTokenClaims, ok := authToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.New("error reading jwt claims")
	}
	if !authToken.Valid {
		return nil, errors.New("auth token is not valid")
	}
	var (
		user   = new(models.User)
		result = s.db.Where("hash =?", authTokenClaims.Subject).Find(user)
	)
	if err := result.Error; err != nil {
		return nil, fmt.Errorf("invalid user")
	}
	return user, nil
}
