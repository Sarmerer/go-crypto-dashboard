package sqlite3

import (
	"log"
	"os"
	"path"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/database"
	"github.com/sarmerer/go-crypto-dashboard/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(portfolio *model.Portfolio) (database.Repository, error) {
	silentLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)

	DBName := config.DefaultDBName
	DBPath := config.DefaultDBPath

	dialector := sqlite.Open(path.Join(DBPath, DBName))

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: silentLogger,
	})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&model.Position{},
		&model.Portfolio{},
		&model.Order{},
		&model.Income{},
	)

	return &repository{db}, nil
}

func (r *repository) limitedDB() *gorm.DB {
	return r.db.Limit(100)
}
