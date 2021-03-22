package user

import (
	"github.com/homepi/homepi/api/services/auth"
	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	auth *auth.Service
}

func NewUserService(db *gorm.DB, authService *auth.Service) *Service {
	return &Service{
		db:   db,
		auth: authService,
	}
}
