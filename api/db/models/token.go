package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/homepi/homepi/api/app/components/strings"
	"gorm.io/gorm"
)

type RefreshedToken struct {
	ID        uint32    `gorm:"primary_key" json:"id,omitempty"`
	TokenId   string    `json:"token_id,omitempty" gorm:"uniqueIndex"`
	Valid     bool      `json:"valid,omitempty" gorm:"default:true"`
	UserId    uint32    `json:"-"`
	User      User      `gorm:"foreignkey:UserId" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// creating random records for user like: stream_key, 2fa_auth_token, etc...
func (user *RefreshedToken) BeforeCreate(scope *gorm.DB) (err error) {
	user.ID = uuid.New().ID()

	// Generating random secrets for user
	user.TokenId = strings.RandomDigits(30)
	user.CreatedAt = time.Now()
	return
}
