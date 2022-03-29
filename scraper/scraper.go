package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/database"
	"github.com/sarmerer/go-crypto-dashboard/driver"
	"github.com/sarmerer/go-crypto-dashboard/model"
	"github.com/sarmerer/go-crypto-dashboard/scraper/exchange"
)

type Scraper interface {
	UpdateExchange(portfolio *model.Portfolio) error

	Scrape() error
	ContinuousScrape()
	ScrapePortfolio(portfolio *model.Portfolio) error
	ScrapePositions() error
	ScrapeOrders() error
	ScrapeIncome() error
	ScrapeDailyBalance() error

	IsWeightOverused() bool
	WaitWeightCooldown()
	Sleep(d time.Duration)
}

type scraper struct {
	repo     driver.Repository
	exchange exchange.Exchange
	ctx      *model.ScrapeCtx
}

func NewScraper() (Scraper, error) {
	return &scraper{
		ctx: &model.ScrapeCtx{
			ScrapedAt:   time.Now().UnixMilli(),
			WeightLimit: 1000,
			Cooldown:    60,
		},
	}, nil
}

func (s *scraper) UpdateExchange(portfolio *model.Portfolio) error {
	exchange, err := exchange.NewExchange(s.ctx.Portfolio, s.ctx)
	if err != nil {
		return err
	}

	s.exchange = exchange
	return nil
}

func (s *scraper) Scrape() (err error) {
	if s.repo, err = database.NewRepository(); err != nil {
		return err
	}

	portfolios, err := config.GetPortfolios()
	if err != nil {
		return err
	}

	for _, portfolio := range portfolios {
		if err := s.ScrapePortfolio(portfolio); err != nil {
			log.Print(fmt.Errorf("portfolio %s: %v", portfolio.Alias, err))
		}
	}

	return nil
}

func (s *scraper) ContinuousScrape() {
	log.Println("continuous scraping started")

	for {
		if err := s.Scrape(); err != nil {
			log.Println("failed to scrape:", err)
		}

		d := time.Minute * 5
		log.Printf("sleeping for %v", d)
		s.divider()
		s.Sleep(d)
	}
}

func (s *scraper) ScrapePortfolio(portfolio *model.Portfolio) error {
	err := s.repo.SyncPortfolio(portfolio)
	if err != nil {
		return err
	}

	s.ctx.Portfolio = portfolio
	log.Printf("scraping portfolio: \"%s\"", s.ctx.Portfolio.ID)

	if err := s.UpdateExchange(portfolio); err != nil {
		return err
	}

	if err := s.repo.RemoveAllOrders(portfolio); err != nil {
		return err
	}

	tasks := []func() error{
		s.ScrapeDailyBalance,
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

func (s *scraper) ScrapeDailyBalance() error {
	log.Println("scraping daily balance")

	balance, err := s.exchange.GetBalance()
	if err != nil {
		return err
	}

	date := time.Now().UTC().Truncate(time.Hour * 24)
	dailyBalance := &model.DailyBalance{
		Balance: balance,
		Date:    date,
	}

	dailyBalance.ScrapeCtx.Apply(s.ctx)
	if err := s.repo.CreateDailyBalance(dailyBalance); err != nil {
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
