package scraper

import (
	"fmt"
	"path"
	"scraper/scraper/exchange"
	"scraper/scraper/models"
	"time"

	"github.com/adshao/go-binance/v2/futures"
	"gorm.io/driver/sqlite"
)

type ScraperConfig struct {
	Portfolios []*models.Portfolio

	DBPath string
}

type Scraper interface {
	shouldUpdateRepo(newPortfolio *models.Portfolio) bool
	shouldUpdateExchange(newPortfolio *models.Portfolio) bool
	UpdateRepo(portfolio *models.Portfolio) error
	UpdateExchange(portfolio *models.Portfolio) error

	Scrape() error
	ScrapePortfolio(portfolio *models.Portfolio) error
	ScrapePositions() error
	ScrapeOrders() error
	ScrapeTrades() error
}

type scraper struct {
	config *ScraperConfig

	repo     Repository
	exchange exchange.Exchange
	ctx      models.ScrapeCtx
	tickers  map[string]*futures.BookTicker
}

func NewScraper(config *ScraperConfig) (Scraper, error) {
	return &scraper{
		config:  config,
		tickers: make(map[string]*futures.BookTicker),
	}, nil
}

func (s *scraper) shouldUpdateRepo(newPortfolio *models.Portfolio) bool {
	oldPath := s.config.DBPath
	if s.ctx.Portfolio.DBPath != "" {
		oldPath = s.ctx.Portfolio.DBPath
	}
	return s.repo == nil || oldPath != newPortfolio.DBPath
}

func (s *scraper) shouldUpdateExchange(newPortfolio *models.Portfolio) bool {
	exchangeMatch := s.ctx.Portfolio.Exchange != newPortfolio.Exchange
	apiKeyMatch := s.ctx.Portfolio.APIKey != newPortfolio.APIKey
	apiSecretMatch := s.ctx.Portfolio.APISecret != newPortfolio.APISecret
	return s.exchange == nil || exchangeMatch || apiKeyMatch || apiSecretMatch
}

func (s *scraper) UpdateRepo(portfolio *models.Portfolio) error {
	shouldUpdate := s.shouldUpdateRepo(portfolio)
	if !shouldUpdate {
		return nil
	}

	DBPath := s.config.DBPath
	if s.ctx.Portfolio.DBPath != "" {
		DBPath = s.ctx.Portfolio.DBPath
	}

	path := path.Join(DBPath, "scraper.db")
	repo, err := NewRepository(sqlite.Open(path))
	if err != nil {
		return err
	}

	s.repo = repo

	portfolios := s.config.Portfolios
	if s.ctx.Portfolio.DBPath != "" {
		portfolios = []*models.Portfolio{s.ctx.Portfolio}
	}

	if err := s.repo.SyncPortfolios(portfolios); err != nil {
		return err
	}

	return nil
}

func (s *scraper) UpdateExchange(portfolio *models.Portfolio) error {
	shouldUpdate := s.shouldUpdateExchange(portfolio)
	if !shouldUpdate {
		return nil
	}

	exchange, err := exchange.NewExchange(s.ctx.Portfolio)
	if err != nil {
		return err
	}

	s.exchange = exchange
	return nil
}

func (s *scraper) Scrape() error {
	s.ctx.ScrapedAt = time.Now().Unix()
	for _, portfolio := range s.config.Portfolios {
		if err := s.ScrapePortfolio(portfolio); err != nil {
			return err
		}
	}

	return nil
}

func (s *scraper) ScrapePortfolio(portfolio *models.Portfolio) error {
	s.ctx.Portfolio = portfolio

	if err := s.UpdateRepo(portfolio); err != nil {
		return err
	}

	if err := s.UpdateExchange(portfolio); err != nil {
		return err
	}

	if err := s.ScrapePositions(); err != nil {
		return err
	}

	if err := s.ScrapeOrders(); err != nil {
		return err
	}

	return nil
}

func (s *scraper) ScrapePositions() error {
	positions, err := s.exchange.GetPositions()
	if err != nil {
		return err
	}

	for _, position := range positions {
		position.ApplyScrapeCtx(&s.ctx)
		if _, err := s.repo.CreatePosition(position); err != nil {
			return fmt.Errorf("failed to create position: %v", err)
		}
	}

	return nil
}

func (s *scraper) ScrapeOrders() error {
	orders, err := s.exchange.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		fmt.Printf("%+v\n", order)
	}

	return nil
}

func (s *scraper) ScrapeTrades() error {
	return nil
}
