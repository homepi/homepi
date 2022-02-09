package models

import (
	"errors"

	"gorm.io/gorm"
)

type Pin struct {
	ID      uint32 `json:"id,omitempty"`
	Disable bool   `json:"disable,omitempty"`
	Used    bool   `json:"used,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Pins struct {
	Top    []*Pin `json:"top,omitempty"`
	Bottom []*Pin `json:"bottom,omitempty"`
}

func (pin *Pin) Check(db *gorm.DB) (used bool, err error) {
	if pin.ID == 0 {
		return false, nil
	}
	var accCount int64
	if err := db.Model(&Accessory{}).Where("pin_id =?", pin.ID).Count(&accCount).Error; err != nil {
		return false, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if accCount == 0 {
		return false, nil
	}
	return true, nil
}
