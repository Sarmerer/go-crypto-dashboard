package database

import (
	"fmt"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/driver"
	"github.com/sarmerer/go-crypto-dashboard/driver/sqlite3"
	"gorm.io/driver/sqlite"
)

func NewRepository() (driver.Repository, error) {
	switch config.DBDriver {
	case "sqlite3":
		dsn := config.SQLite3Config.DSN
		dialector := sqlite.Open(dsn)
		return sqlite3.NewRepository(dialector)
	default:
		return nil, fmt.Errorf("unknown database driver: %s", config.DBDriver)
	}
}
