package exchange

import (
	"fmt"

	"github.com/sarmerer/go-crypto-dashboard/tracker/model"
)

type Exchange interface {
	GetSymbolPrices() ([]*model.SymbolPrice, error)
	GetBalance() (float64, error)
	GetPositions() ([]*model.Position, error)
	GetOrders() ([]*model.Order, error)
	GetIncome() ([]*model.Income, error)
	GetIncomeBetween(start, end int64) ([]*model.Income, error)
}

func NewExchange(portfolio *model.Portfolio, ctx *model.ScrapeCtx) (Exchange, error) {
	switch portfolio.Exchange {
	case "binance-futures":
		return NewBinanceFutures(portfolio, ctx)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", portfolio.Exchange)
	}
}
