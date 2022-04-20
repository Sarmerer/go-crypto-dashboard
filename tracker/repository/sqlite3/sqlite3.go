package sqlite3

import (
	"log"
	"os"

	"github.com/sarmerer/go-crypto-dashboard/tracker/model"
	"github.com/sarmerer/go-crypto-dashboard/tracker/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DSN string
}

type repo struct {
	db *gorm.DB
}

func NewRepository(dialector gorm.Dialector) (repository.Repository, error) {
	silentLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			IgnoreRecordNotFoundError: true,
		},
	)

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
		&model.DailyBalance{},
		&model.CurrentBalance{},
		&model.SymbolPrice{},
	)

	return &repo{db}, nil
}

func (r *repo) limitedDB() *gorm.DB {
	return r.db.Limit(100)
}
