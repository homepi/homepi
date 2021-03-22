package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID                 uint32    `gorm:"primary_key" json:"id"`
	Title              string    `json:"title" gorm:"uniqueIndex"`
	Administrator      bool      `json:"administrator" gorm:"default:false"`
	CanSeeAccessories  bool      `json:"can_see_accessories" gorm:"default:false"`
	CanRunAccessory    bool      `json:"can_run_accessory" gorm:"default:false"`
	CanCreateAccessory bool      `json:"can_create_accessory" gorm:"default:false"`
	CanRemoveAccessory bool      `json:"can_remove_accessory" gorm:"default:false"`
	CanSeeWebhook      bool      `json:"can_see_webhook" gorm:"default:false"`
	CanRemoveWebhook   bool      `json:"can_remove_webhook" gorm:"default:false"`
	CanCreateWebhook   bool      `json:"can_create_webhook" gorm:"default:false"`
	CanSeeUsers        bool      `json:"can_see_users" gorm:"default:false"`
	CanCreateUser      bool      `json:"can_create_user" gorm:"default:false"`
	CanRemoveUser      bool      `json:"can_remove_user" gorm:"default:false"`
	CanSeeRoles        bool      `json:"can_see_roles" gorm:"default:false"`
	CanCreateRole      bool      `json:"can_create_role" gorm:"default:false"`
	CanRemoveRole      bool      `json:"can_remove_role" gorm:"default:false"`
	CanSeeLogs         bool      `json:"can_see_logs" gorm:"default:false"`
	CreatedAt          time.Time `json:"-"`
	UpdatedAt          time.Time `json:"-"`
}

func GetRoleByName(name string) *Role {
	switch name {
	case "admin":
		return AdminRole(name)
	default:
		return UserRole(name)
	}
}

func AdminRole(title string) *Role {
	return &Role{
		Title:         title,
		Administrator: true,
	}
}

func UserRole(title string) *Role {
	return &Role{
		Title:             title,
		CanSeeAccessories: true,
	}
}

func (r *Role) BeforeCreate(scope *gorm.DB) (err error) {
	r.ID = uuid.New().ID()
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return
}
