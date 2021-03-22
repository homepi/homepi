package accessory

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewAccessoryService(db *gorm.DB) *Service {
	return &Service{db: db}
}
