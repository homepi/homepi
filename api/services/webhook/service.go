package webhook

import "gorm.io/gorm"

type Service struct {
	db *gorm.DB
}

func NewWebhookService(db *gorm.DB) *Service {
	return &Service{db: db}
}
