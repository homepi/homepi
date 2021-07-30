package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"gorm.io/gorm"
)

type APIToken struct {
	ID        uint32    `gorm:"primary_key" json:"id"`
	Token     string    `json:"token,omitempty" gorm:"uniqueIndex"`
	UserID    uint32    `json:"-"`
	User      *User     `gorm:"foreignkey:UserID" json:"-" validate:"-"`
	RoleID    uint32    `json:"-"`
	Role      *Role     `gorm:"foreignkey:RoleID" json:"-" validate:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (t *APIToken) BeforeCreate(db *gorm.DB) error {
	t.ID = uuid.New().ID()
	t.Token = libstr.RandomLetters(60)
	return nil
}
