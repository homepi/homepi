package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        int64     `gorm:"primary_key" json:"id"`
	Fullname  string    `json:"fullname"  form:"fullname"  validate:"required"`
	Username  string    `json:"username"  form:"username"  validate:"required" gorm:"uniqueIndex"`
	Email     string    `json:"email"     form:"email"     validate:"required,email" gorm:"uniqueIndex"`
	Password  string    `json:"-"         form:"password"  validate:"required"`
	Avatar    string    `json:"avatar"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	RoleID    uint32    `json:"-"`
	Role      *Role     `gorm:"foreignkey:RoleID" json:"role"`
	LastLogin time.Time `json:"last_login"`
	JoinedAt  time.Time `json:"joined_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchUser struct {
	Query string `json:"query" form:"query" validate:"required"`
}

// Validate users's password
func (user *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// Generate random users hash id
func (user *User) GenerateUserHash() string {
	return libstr.RandomLetters(30)
}

// Hash users's password with bcrypt
func (User) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// set users password
func (user *User) SetPassword(password string) (err error) {
	user.Password, err = user.HashPassword(password)
	return
}

func (user *User) LoadRelations(db *gorm.DB) {
	db.Model(user).Association("Role")
}

// creating random records for user like: stream_key, 2fa_auth_token, etc...
func (user *User) BeforeCreate(db *gorm.DB) error {
	user.ID = int64(uuid.New().ID())
	user.Avatar = "default"
	user.JoinedAt = time.Now()
	user.LastLogin = time.Now()
	return nil
}
