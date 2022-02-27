package scraper

import (
	"fmt"
	"reflect"
	"scraper/scraper/models"

	"gorm.io/gorm"
)

type Repository interface {
	SyncPortfolios(portfolios []*models.Portfolio) error
	SyncPortfolio(portfolio *models.Portfolio) error
	CreatePosition(position *models.Position) (*models.Position, error)

	UpdateOrCreate(object interface{}) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(dialector gorm.Dialector) (Repository, error) {
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&models.Position{},
		&models.Portfolio{},
		&models.Order{},
		&models.Trade{},
	)

	return &repository{db}, nil
}

func (r *repository) SyncPortfolios(portfolios []*models.Portfolio) error {
	for _, portfolio := range portfolios {
		if err := r.SyncPortfolio(portfolio); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) SyncPortfolio(portfolio *models.Portfolio) error {
	if err := r.UpdateOrCreate(portfolio); err != nil {
		return err
	}

	return nil
}

func (r *repository) CreatePosition(position *models.Position) (*models.Position, error) {
	if err := r.UpdateOrCreate(position); err != nil {
		return nil, err
	}
	return position, nil
}

func (r *repository) UpdateOrCreate(object interface{}) error {
	res := r.db.Updates(object)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		if err := r.db.Create(object).Error; err != nil {
			objectType := reflect.TypeOf(object).String()
			return fmt.Errorf("failed to create %s: %v", objectType, err)
		}
	}

	return nil
}
