package scraper

import "github.com/adshao/go-binance/v2/futures"

type ScraperConfig struct {
	APIKey    string
	APISecret string

	DBPath string
}

type Scrape struct {
	Timestamp int64
	Account   string
}

type Scraper interface {
	Scrape() error

	GetPositions() ([]*futures.Balance, error)
	GetOrders() ([]*futures.Order, error)
	GetTrades(symbol string) ([]*futures.Trade, error)
	GetPositionCost(Position *futures.Balance) (float64, error)

	AssetIsUSD(asset string) bool

	updateTickers() error
}

type Position struct {
	Scrape

	Symbol string `gorm:"primary_key"`
	Amount float64
	Price  float64
}

type Order struct {
	Scrape

	Symbol string
	Price  float64
	Amount float64
}

type Trade struct {
	Scrape

	Symbol string
	Price  float64
	Amount float64
}
