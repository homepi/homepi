package system

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) *Service {
	return &Service{db: db}
}
