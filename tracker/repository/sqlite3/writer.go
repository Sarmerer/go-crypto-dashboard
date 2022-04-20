package sqlite3

import (
	"errors"
	"reflect"

	"github.com/sarmerer/go-crypto-dashboard/tracker/model"
	"gorm.io/gorm"
)

func (r *repo) CreateSymbolPrice(sp *model.SymbolPrice) error {
	return r.createOrUpdate(sp, "symbol = ?", sp.Symbol)
}

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
	return r.createOrUpdate(position, "symbol = ? AND side = ? AND portfolio_id = ?", position.Symbol, position.Side, position.Portfolio.ID)
}

func (r *repo) CreateOrder(order *model.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return err
	}

	return nil
}

func (r *repo) RemoveAllPositions(portfolio *model.Portfolio) error {
	return r.db.Where("portfolio_id = ?", portfolio.ID).Delete(model.Position{}).Error
}

func (r *repo) RemoveAllOrders(portfolio *model.Portfolio) error {
	return r.db.Where("portfolio_id = ?", portfolio.ID).Delete(&model.Order{}).Error
}

func (r *repo) CreateIncome(income *model.Income) error {
	return r.createOrUpdate(income, "id = ? AND portfolio_id = ? AND type = ?", income.ID, income.Portfolio.ID, income.Type)
}

func (r *repo) CreateDailyBalance(balance *model.DailyBalance) error {
	return r.createOrUpdate(balance, "date = ? AND portfolio_id = ?", balance.Date, balance.Portfolio.ID)
}

func (r *repo) UpdateCurrentBalance(balance *model.CurrentBalance) error {
	return r.createOrUpdate(balance, "portfolio_id = ?", balance.Portfolio.ID)
}

func (r *repo) createOrUpdate(model interface{}, query string, args ...interface{}) error {
	mt := reflect.TypeOf(model)
	dummy := reflect.New(mt).Interface()
	if err := r.db.Where(query, args...).First(dummy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.Create(model).Error
		}
		return err
	}

	return r.db.Model(dummy).Updates(model).Error
}
