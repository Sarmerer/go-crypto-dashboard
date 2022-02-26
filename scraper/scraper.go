package scraper

import (
	"context"
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type scraper struct {
	db      *gorm.DB
	client  *futures.Client
	tickers map[string]*futures.BookTicker
}

func New(config *ScraperConfig) (Scraper, error) {

	db, err := initDB(config)
	if err != nil {
		return nil, err
	}

	client, err := initClient(config)
	if err != nil {
		return nil, err
	}

	return &scraper{
		db:      db,
		client:  client,
		tickers: make(map[string]*futures.BookTicker),
	}, nil
}

func initDB(config *ScraperConfig) (*gorm.DB, error) {
	path := path.Join(config.DBPath, "scraper.db")
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Position{}, &Trade{}, &Order{})

	return db, nil
}

func initClient(config *ScraperConfig) (*futures.Client, error) {
	client := binance.NewFuturesClient(config.APIKey, config.APISecret)
	err := client.NewPingService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping binance api: %v", err)
	}

	return client, nil
}

func (s *scraper) Scrape() error {
	err := s.updateTickers()
	if err != nil {
		return err
	}

	positions, err := s.GetPositions()
	if err != nil {
		return err
	}

	for _, position := range positions {
		err = s.db.Create(position).Error
		if err != nil {
			return err
		}

		// pair := position.Asset + "USDT"
		// trades, err := s.GetTrades(pair)
		// if err != nil {
		// 	return err
		// }

		// for _, trade := range trades {
		// 	err = s.db.Create(&Trade{
		// 		Symbol:    pair,
		// 		Price:     trade.Price,
		// 		Amount:    trade.Quantity,
		// 		Timestamp: trade.Time,
		// 		Side:      trade.IsBuyer,
		// 	}).Error
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}

	return nil
}

func (s *scraper) updateTickers() error {
	tickers, err := s.client.NewListBookTickersService().Do(context.Background())
	if err != nil {
		return err
	}

	s.tickers = make(map[string]*futures.BookTicker)
	for _, ticker := range tickers {
		s.tickers[ticker.Symbol] = ticker
	}

	return nil
}

func (s *scraper) GetPositions() ([]*futures.Balance, error) {
	balance, err := s.client.NewGetBalanceService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	positions := []*futures.Balance{}
	for _, asset := range balance {

		if s.AssetIsUSD(asset.Asset) {
			continue
		}

		pair := asset.Asset + "USDT"
		ticker := s.tickers[pair]
		if ticker == nil {
			return nil, fmt.Errorf("ticker not found for %s", pair)
		}

		positionCost, err := s.GetPositionCost(asset)
		if err != nil {
			return nil, err
		}

		if positionCost < 0.1 {
			continue
		}

		positions = append(positions, asset)
	}
	return positions, nil
}

func (s *scraper) GetOrders() ([]*futures.Order, error) {
	return s.client.NewListOpenOrdersService().Do(context.Background())
}

func (s *scraper) GetTrades(symbol string) ([]*futures.Trade, error) {
	return s.client.NewRecentTradesService().Symbol(symbol).Do(context.Background())
}

func (s *scraper) GetPositionCost(position *futures.Balance) (float64, error) {

	pair := position.Asset + "USDT"
	ticker := s.tickers[pair]
	if ticker == nil {
		return 0, fmt.Errorf("ticker not found for %s", pair)
	}

	b, err := strconv.ParseFloat(position.Balance, 64)
	if err != nil {
		return 0, err
	}

	p, err := strconv.ParseFloat(ticker.BidPrice, 64)
	if err != nil {
		return 0, err
	}

	return b * p, nil
}

func (s *scraper) AssetIsUSD(asset string) bool {
	fiatUSD := []string{
		"USDT",
		"USDC",
		"TUSD",
		"BUSD",
	}

	return contains(fiatUSD, asset)
}

func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
