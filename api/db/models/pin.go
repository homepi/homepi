package models

import (
	"errors"

	"gorm.io/gorm"
)

type Pin struct {
	Id      uint32 `json:"id,omitempty"`
	Disable bool   `json:"disable,omitempty"`
	Used    bool   `json:"used,omitempty"`
	Name    string `json:"name,omitempty"`
}

type Pins struct {
	Top    []*Pin `json:"top,omitempty"`
	Bottom []*Pin `json:"bottom,omitempty"`
}

func (pin *Pin) Check(db *gorm.DB) (used bool, err error) {
	accs := new(Accessory)
	result := db.Where("pin_id =?", pin.Id).Find(accs)
	if err := result.Error; err != nil {
		return false, err
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, nil
}
