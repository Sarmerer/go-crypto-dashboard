package sqlite3

import (
	"errors"

	"github.com/sarmerer/go-crypto-dashboard/model"
	"gorm.io/gorm"
)

func (r *repo) SyncPortfolio(portfolio *model.Portfolio) error {
	record := &model.Portfolio{}
	err := r.db.Where("id = ?", portfolio.ID).First(record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(portfolio).Error
		}
		return err
	}

	portfolio.SyncWith(record)
	return r.db.Model(portfolio).Updates(portfolio).Error
}

func (r *repo) UpdatePortfolio(portfolio *model.Portfolio) error {
	return r.db.Model(portfolio).Updates(portfolio).Error
}

func (r *repo) CreatePosition(position *model.Position) error {
	return r.createIfNotExists(position, "symbol = ? AND side = ? AND portfolio_id = ?", position.Symbol, position.Side, position.Portfolio.ID)
}

func (r *repo) CreateOrder(order *model.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func (r *repo) RemoveAllOrders(portfolio *model.Portfolio) error {
	return r.db.Where("portfolio_id = ?", portfolio.ID).Delete(&model.Order{}).Error
}

func (r *repo) CreateIncome(income *model.Income) error {
	return r.createIfNotExists(income, "id = ? AND portfolio_id = ?", income.ID, income.Portfolio.ID)
}

func (r *repo) CreateDailyBalance(balance *model.DailyBalance) error {
	return r.createIfNotExists(balance, "date = ? AND portfolio_id = ?", balance.Date, balance.Portfolio.ID)
}

func (r *repo) createIfNotExists(model interface{}, query string, args ...interface{}) error {
	err := r.db.Model(model).Where(query, args...).First(model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(model).Error
	}

	return err
}
