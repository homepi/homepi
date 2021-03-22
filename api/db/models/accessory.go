package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/stianeikeland/go-rpio"
	"gorm.io/gorm"
)

type Accessory struct {
	ID          uint32     `gorm:"primary_key" json:"id"`
	Name        string     `json:"name,omitempty" form:"name" validate:"required"`
	Task        Task       `json:"task,omitempty" form:"task" validate:"accessory_task"`
	Description string     `json:"description,omitempty" form:"description" validate:"required"`
	Icon        string     `json:"icon,omitempty"`
	PinId       uint64     `json:"pin_id,omitempty" form:"pin_id" validate:"required" gorm:"uniqueIndex"`
	IsActive    bool       `json:"is_active,omitempty"`
	IsPublic    bool       `json:"is_public,omitempty" form:"is_public" validate:"required"`
	UserId      uint32     `json:"-"`
	User        User       `gorm:"foreignkey:UserId" json:"-" validate:"-"`
	State       rpio.State `json:"state,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty"`
}

type AccessoryWithUserData struct {
	*Accessory
	User User `gorm:"foreignkey:UserId" json:"user,omitempty" validate:"-"`
}

func (ac *Accessory) UpdateStatus(db *gorm.DB, status rpio.State) {
	db.Model(ac).Update("status", status)
}

func (ac *Accessory) GpioPin() (pin rpio.Pin) {
	pin = rpio.Pin(ac.PinId)
	pin.Output()
	pin.Clock()
	return
}

func (ac *AccessoryWithUserData) LoadRelations(db *gorm.DB) {
	db.Model(ac).Association("User")
	db.Model(ac.User).Association("Role")
}

func (ac *Accessory) BeforeCreate(db *gorm.DB) (err error) {
	ac.ID = uuid.New().ID()
	ac.CreatedAt = time.Now()
	ac.UpdatedAt = time.Now()
	return
}
