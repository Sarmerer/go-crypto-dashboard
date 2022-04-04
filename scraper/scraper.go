package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/model"
	"github.com/sarmerer/go-crypto-dashboard/repository"
	"github.com/sarmerer/go-crypto-dashboard/scraper/exchange"
)

type Scraper interface {
	GetExchange(portfolio *model.Portfolio) (exchange.Exchange, error)

	Scrape() error
	ContinuousScrape() error
	ScrapePrices(portfolio *model.Portfolio) error
	ScrapePortfolio(portfolio *model.Portfolio) error
	ScrapePositions() error
	ScrapeOrders() error
	ScrapeIncome() error
	ScrapeBalance() error

	IsWeightOverused() bool
	WaitWeightCooldown()
	Sleep(d time.Duration)
}

type scraper struct {
	repo     repository.Repository
	exchange exchange.Exchange
	ctx      *model.ScrapeCtx
}

func NewScraper(repo repository.Repository) (Scraper, error) {
	return &scraper{
		ctx: &model.ScrapeCtx{
			ScrapedAt:   time.Now().UnixMilli(),
			WeightLimit: 1000,
			Cooldown:    60,
		},
		repo: repo,
	}, nil
}

func (s *scraper) GetExchange(portfolio *model.Portfolio) (exchange.Exchange, error) {
	exchange, err := exchange.NewExchange(portfolio, s.ctx)
	if err != nil {
		return nil, err
	}

	return exchange, nil
}

func (s *scraper) Scrape() (err error) {
	portfolios, err := config.GetPortfolios()
	if err != nil {
		return err
	}

	if err := s.ScrapePrices(portfolios[0]); err != nil {
		return err
	}

	for _, portfolio := range portfolios {
		if err := s.ScrapePortfolio(portfolio); err != nil {
			log.Print(fmt.Errorf("portfolio %s: %v", portfolio.Alias, err))
		}
	}

	return nil
}

func (s *scraper) ContinuousScrape() error {
	log.Println("continuous scraping started")

	if err := s.Scrape(); err != nil {
		return err
	}

	for {
		d := time.Minute * 5
		log.Printf("sleeping for %v", d)
		s.divider()
		s.Sleep(d)

		if err := s.Scrape(); err != nil {
			log.Println("failed to scrape:", err)
		}
	}
}

func (s *scraper) ScrapePortfolio(portfolio *model.Portfolio) error {
	err := s.repo.SyncPortfolio(portfolio)
	if err != nil {
		return err
	}

	s.ctx.Portfolio = portfolio
	log.Printf("scraping portfolio: \"%s\"", s.ctx.Portfolio.ID)

	exchange, err := s.GetExchange(portfolio)
	if err != nil {
		return err
	}
	s.exchange = exchange

	if err := s.repo.RemoveAllPositions(portfolio); err != nil {
		return err
	}

	if err := s.repo.RemoveAllOrders(portfolio); err != nil {
		return err
	}

	tasks := []func() error{
		s.ScrapeBalance,
		s.ScrapePositions,
		s.ScrapeOrders,
		s.ScrapeIncome,
	}

	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}

		if s.IsWeightOverused() {
			s.WaitWeightCooldown()
		}
	}

	s.divider()
	return nil
}

func (s *scraper) ScrapePrices(portfolio *model.Portfolio) error {
	log.Printf("scraping prices from %s\n", portfolio.Exchange)
	s.divider()

	exchange, err := s.GetExchange(portfolio)
	if err != nil {
		return err
	}

	prices, err := exchange.GetSymbolPrices()
	if err != nil {
		return err
	}

	for _, price := range prices {
		price.ScrapeCtx.Apply(s.ctx)
		if err := s.repo.CreateSymbolPrice(price); err != nil {
			return err
		}
	}

	return nil
}

func (s *scraper) ScrapePositions() error {
	log.Println("scraping positions")

	positions, err := s.exchange.GetPositions()
	if err != nil {
		return err
	}

	for _, position := range positions {
		position.ScrapeCtx.Apply(s.ctx)
		if err := s.repo.CreatePosition(position); err != nil {
			return err
		}

	}

	return nil
}

func (s *scraper) ScrapeOrders() error {
	log.Println("scraping orders")

	orders, err := s.exchange.GetOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		order.ScrapeCtx.Apply(s.ctx)
		if err := s.repo.CreateOrder(order); err != nil {
			return err
		}
	}

	return nil
}

func (s *scraper) ScrapeIncome() error {
	if !s.ctx.Portfolio.HistoryScraped {
		return s.scrapeIncomeHistory()
	}

	log.Println("scraping recent income")

	incomes, err := s.exchange.GetIncome()
	if err != nil {
		return err
	}

	for _, income := range incomes {
		income.ScrapeCtx.Apply(s.ctx)
		if err := s.repo.CreateIncome(income); err != nil {
			return err
		}
	}

	return nil
}

func (s *scraper) scrapeIncomeHistory() error {
	log.Println("scraping historical income")

	oldestIncomeTime := time.Now().UnixMilli()
	for {
		if s.IsWeightOverused() {
			s.WaitWeightCooldown()
			log.Println("scraping next chunk...")
		}

		incomes, err := s.exchange.GetIncomeBetween(0, oldestIncomeTime)
		if err != nil {
			return err
		}

		if len(incomes) == 0 {
			break
		}

		for _, income := range incomes {
			income.ScrapeCtx.Apply(s.ctx)
			if err := s.repo.CreateIncome(income); err != nil {
				return err
			}
		}

		newOldest := incomes[0].Date.UnixMilli()
		if newOldest >= oldestIncomeTime {
			break
		}

		oldestIncomeTime = newOldest - 1
	}

	s.ctx.Portfolio.HistoryScraped = true
	if err := s.repo.UpdatePortfolio(s.ctx.Portfolio); err != nil {
		return err
	}

	return nil
}

func (s *scraper) ScrapeBalance() error {
	log.Println("scraping balance")

	date := time.Now().UTC()
	balance, err := s.exchange.GetBalance()
	if err != nil {
		return err
	}

	dailyBalance := &model.DailyBalance{
		Balance: balance,
		Date:    date.Truncate(time.Hour * 24),
	}

	dailyBalance.ScrapeCtx.Apply(s.ctx)
	if err := s.repo.CreateDailyBalance(dailyBalance); err != nil {
		return err
	}

	currentBalance := &model.CurrentBalance{Balance: balance, Date: date}
	currentBalance.ScrapeCtx.Apply(s.ctx)
	if err := s.repo.UpdateCurrentBalance(currentBalance); err != nil {
		return err
	}

	return nil
}

func (s *scraper) IsWeightOverused() bool {
	return s.ctx.WeightUsed > s.ctx.WeightLimit
}

func (s *scraper) WaitWeightCooldown() {
	log.Printf("used weight: %d/%d, sleeping for %d seconds", s.ctx.WeightUsed, s.ctx.WeightLimit, s.ctx.Cooldown)
	s.Sleep(time.Second * s.ctx.Cooldown)
	s.ctx.WeightUsed = 0
}

func (s *scraper) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (s *scraper) divider() {
	log.Println("------------------")
}
