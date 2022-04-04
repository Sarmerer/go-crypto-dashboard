package repository

import (
	"github.com/sarmerer/go-crypto-dashboard/model"
)

type Repository interface {
	Reader
	Writer
}

type Reader interface {
	GetPositions() ([]*model.Position, error)
	GetPortfolios() ([]*model.Portfolio, error)
	GetOrders() ([]*model.Order, error)
	GetIncomeBetween(start, end int64) ([]*model.Income, error)
}

type Writer interface {
	SyncPortfolio(portfolio *model.Portfolio) error
	UpdatePortfolio(portfolio *model.Portfolio) error

	CreateSymbolPrice(price *model.SymbolPrice) error
	CreatePosition(position *model.Position) error
	CreateOrder(order *model.Order) error
	CreateIncome(income *model.Income) error
	CreateDailyBalance(balance *model.DailyBalance) error
	UpdateCurrentBalance(balance *model.CurrentBalance) error

	RemoveAllPositions(portfolio *model.Portfolio) error
	RemoveAllOrders(portfolio *model.Portfolio) error
}
