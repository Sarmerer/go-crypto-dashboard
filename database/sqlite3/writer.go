package sqlite3

import (
	"errors"

	"github.com/sarmerer/go-crypto-dashboard/model"
	"gorm.io/gorm"
)

func (r *repository) UpdatePortfolios(portfolios []*model.Portfolio) error {
	for _, portfolio := range portfolios {
		if err := r.UpdatePortfolio(portfolio); err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) UpdatePortfolio(portfolio *model.Portfolio) error {
	err := r.db.Where("id = ?", portfolio.ID).First(portfolio).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(portfolio).Error
	}

	if err != nil {
		return err
	}

	err = r.db.Model(portfolio).Updates(portfolio).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreatePosition(position *model.Position) error {
	return r.createIfNotExists(position, "symbol = ? AND side = ?", position.Symbol, position.Side)
}

func (r *repository) CreateOrder(order *model.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) RemoveAllOrders(portfolio *model.Portfolio) error {
	return r.db.Where("portfolio_id = ?", portfolio.ID).Delete(&model.Order{}).Error
}

func (r *repository) CreateIncome(income *model.Income) error {
	return r.createIfNotExists(income, "id = ?", income.ID)
}

func (r *repository) createIfNotExists(model interface{}, query string, args ...interface{}) error {
	err := r.db.Model(model).Where(query, args...).First(model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(model).Error
	}

	return err
}
