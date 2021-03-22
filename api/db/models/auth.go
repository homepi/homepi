package models

type Auth struct {
	User     string   `form:"user"  validate:"required"`
	Pass     string   `form:"pass"  validate:"required"`
}