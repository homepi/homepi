package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uint32    `gorm:"primary_key" json:"id,omitempty"`
	TokenID   string    `json:"token_id,omitempty" gorm:"uniqueIndex"`
	Valid     bool      `json:"valid,omitempty" gorm:"default:true"`
	UserID    int64     `json:"-"`
	User      *User     `gorm:"foreignkey:UserID" json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// creating random records for user like: stream_key, 2fa_auth_token, etc...
func (user *RefreshToken) BeforeCreate(scope *gorm.DB) (err error) {
	user.ID = uuid.New().ID()

	// Generating random secrets for user
	user.TokenID = libstr.RandomDigits(30)
	user.CreatedAt = time.Now()
	return
}
