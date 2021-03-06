package sqlite3

import (
	"github.com/sarmerer/go-crypto-dashboard/tracker/model"
)

func (r *repo) GetPortfolios() ([]*model.Portfolio, error) {
	var portfolios []*model.Portfolio
	if err := r.limitedDB().Find(&portfolios).Error; err != nil {
		return nil, err
	}

	return portfolios, nil
}

func (r *repo) GetPositions() ([]*model.Position, error) {
	var positions []*model.Position
	if err := r.limitedDB().Find(&positions).Error; err != nil {
		return nil, err
	}

	return positions, nil
}

func (r *repo) GetOrders() ([]*model.Order, error) {
	var orders []*model.Order
	if err := r.limitedDB().Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repo) GetIncomeBetween(start, end int64) ([]*model.Income, error) {
	var incomes []*model.Income
	err := r.limitedDB().
		Where("(type <> 'TRANSFER' OR type is null) AND  timestamp >= ? AND timestamp <= ?", start, end).
		Find(&incomes).
		Error
	if err != nil {
		return nil, err
	}

	return incomes, nil
}
