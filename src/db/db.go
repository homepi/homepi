package db

import (
	"fmt"

	"github.com/homepi/homepi/src/core"
	"gorm.io/gorm"

	"gorm.io/driver/sqlite"
)

func NewConnection(cfg *core.ConfMap) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(cfg.DB.Path), nil)
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite3 database CLI: %s", err)
	}
	return database, nil
}
