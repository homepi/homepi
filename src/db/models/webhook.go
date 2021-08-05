package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"gorm.io/gorm"
)

type Webhook struct {
	ID          uint32     `gorm:"primary_key" json:"id"`
	Name        string     `json:"name" validate:"required" form:"name"`
	Hash        string     `json:"hash"`
	IsPublic    bool       `json:"is_public" gorm:"default:false" form:"is_public"`
	IsActive    bool       `json:"is_active" gorm:"default:true" form:"is_active"`
	AccessoryID uint32     `json:"-" validate:"required" form:"accessory_id"`
	Accessory   *Accessory `gorm:"foreignkey:AccessoryID" json:"accessory" validate:"-"`
	UserID      uint32     `json:"-"`
	User        *User      `gorm:"foreignkey:UserID" json:"-" validate:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (wh *Webhook) GenerateHash() string {
	return libstr.RandomLetters(25)
}

func (wh *Webhook) LoadRelations(db *gorm.DB) {
	db.Model(wh).Association("User")
	db.Model(wh).Association("Accessory")
}

func GetWebhooks(db *gorm.DB, user *User, limit int) (webhooks []*Webhook, err error) {
	var (
		result = db.Where("user_id =?", user.ID).
			Order("created_at desc").
			Limit(limit).
			Find(&webhooks)
	)
	if err = result.Error; err != nil {
		return
	}
	for _, w := range webhooks {
		w.LoadRelations(db)
	}
	return
}

// creating random records for user like: stream_key, 2fa_auth_token, etc...
func (wh *Webhook) BeforeCreate(db *gorm.DB) (err error) {
	wh.ID = uuid.New().ID()
	wh.CreatedAt = time.Now()
	wh.UpdatedAt = time.Now()
	return
}
