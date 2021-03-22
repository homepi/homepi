package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Log struct {
	ID          uint32    `gorm:"primary_key" json:"id"`
	Read        bool      `json:"read,omitempty"`
	Type        LogType   `json:"type"`
	UserId      uint32    `json:"-"`
	User        User      `gorm:"foreignkey:UserId" json:"user"`
	WebhookId   uint32    `json:"-"`
	Webhook     Webhook   `gorm:"foreignkey:WebhookId" json:"-"`
	AccessoryId uint32    `json:"-"`
	Accessory   Accessory `gorm:"foreignkey:AccessoryId" json:"accessory"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type LogWithUser struct {
	*Log
	User User `gorm:"foreignkey:UserId" json:"user,omitempty" validate:"-"`
}

type LogWithWebhook struct {
	*Log
	Webhook Webhook `gorm:"foreignkey:WebhookId" json:"webhook,omitempty" validate:"-"`
}

type LogType int

const (
	//InvalidLogType   LogType = 0
	UserLogType LogType = 1
	LogWebHook  LogType = 2
)

func (l *Log) LoadRelations(db *gorm.DB) {
	db.Model(l).Association("User")
	db.Model(l).Association("Accessory")
	if l.Type == LogWebHook {
		db.Model(l).Association("Webhook")
		db.Model(l.Webhook).Association("User")
	}
}

func GetLogs(db *gorm.DB, user *User, limit int) (logs []interface{}, err error) {

	var (
		dbLogs = make([]*Log, 0)
		result = db.Where("user_id =?", user.ID).
			Order("created_at desc").
			Limit(limit).
			Find(&dbLogs)
	)

	if err = result.Error; err != nil {
		return
	}

	for _, l := range dbLogs {
		l.LoadRelations(db)
		var resultLog interface{}
		switch l.Type {
		case LogWebHook:
			resultLog = &LogWithWebhook{
				Log:     l,
				Webhook: l.Webhook,
			}
		default:
			resultLog = l
		}
		logs = append(logs, resultLog)
	}

	return
}

func (l *Log) BeforeCreate(scope *gorm.DB) (err error) {
	l.ID = uuid.New().ID()
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()
	return
}
