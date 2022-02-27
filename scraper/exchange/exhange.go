package exchange

import (
	"fmt"
	"scraper/scraper/models"

	"github.com/adshao/go-binance/v2/futures"
)

type Exchange interface {
	GetPositions() ([]*models.Position, error)
	GetOrders() ([]*models.Order, error)
	GetTrades(symbol string) ([]*models.Trade, error)

	parsePosition(accountPosition *futures.AccountPosition) (*models.Position, error)
	parseOrder(order *futures.Order) (*models.Order, error)
	parseTrade(trade *futures.Trade) (*models.Trade, error)
}

func NewExchange(portfolio *models.Portfolio) (Exchange, error) {
	switch portfolio.Exchange {
	case "binance-futures":
		return NewBinanceFutures(portfolio)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", portfolio.Exchange)
	}
}
