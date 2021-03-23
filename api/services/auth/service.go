package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Token struct {
	Secret     []byte
	ExpireTime time.Duration
}

type JWTConfig struct {
	AccessToken  *Token
	RefreshToken *Token
}

type Service struct {
	db  *gorm.DB
	jwt *JWTConfig
}

func NewAuthService(db *gorm.DB) (*Service, error) {

	accessTokenSecret := os.Getenv("HPI_ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		return nil, fmt.Errorf("HPI_ACCESS_TOKEN_SECRET is required!")
	}

	var (
		err                      error
		accessTokenExpireTime    = 240
		accessTokenExpireTimeStr = os.Getenv("HPI_ACCESS_TOKEN_EXPIRE_TIME")
	)

	accessTokenExpireTimeStr = os.Getenv("HPI_ACCESS_TOKEN_EXPIRE_TIME")
	if accessTokenExpireTimeStr != "" {
		accessTokenExpireTime, err = strconv.Atoi(accessTokenExpireTimeStr)
		if err != nil {
			return nil, fmt.Errorf("HPI_ACCESS_TOKEN_EXPIRE_TIME is required or should be a number!")
		}
	}

	refreshTokenSecret := os.Getenv("HPI_REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		return nil, fmt.Errorf("HPI_REFRESH_TOKEN_SECRET is required!")
	}

	var (
		refreshTokenExpireTime    = 1440
		refreshTokenExpireTimeStr = os.Getenv("REFRESH_TOKEN_EXPIRE_TIME")
	)

	if refreshTokenExpireTimeStr != "" {
		refreshTokenExpireTime, err = strconv.Atoi(refreshTokenExpireTimeStr)
		if err != nil {
			return nil, errors.New("cannot convert jwt refresh token expire time type string to int")
		}
	}

	jwtConfig := &JWTConfig{
		RefreshToken: &Token{
			Secret:     []byte(accessTokenSecret),
			ExpireTime: time.Duration(accessTokenExpireTime),
		},
		AccessToken: &Token{
			Secret:     []byte(refreshTokenSecret),
			ExpireTime: time.Duration(refreshTokenExpireTime),
		},
	}

	return &Service{
		db:  db,
		jwt: jwtConfig,
	}, nil
}
