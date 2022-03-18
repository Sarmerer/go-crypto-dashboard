package scraper

import (
	"fmt"
	"log"
	"time"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/database"
	"github.com/sarmerer/go-crypto-dashboard/database/sqlite3"
	"github.com/sarmerer/go-crypto-dashboard/model"
	"github.com/sarmerer/go-crypto-dashboard/scraper/exchange"
)

type Scraper interface {
	UpdateRepo(portfolio *model.Portfolio) error
	UpdateExchange(portfolio *model.Portfolio) error

	Scrape() error
	ScrapePortfolio(portfolio *model.Portfolio) error
	ScrapePositions() error
	ScrapeOrders() error
	ScrapeIncome() error
	ScrapeIncomeHistory() error

	IsWeightOverused() bool
	Sleep()
}

type scraper struct {
	repo     database.Repository
	exchange exchange.Exchange
	ctx      *model.ScrapeCtx
}

func NewScraper() (Scraper, error) {
	return &scraper{
		ctx: &model.ScrapeCtx{
			ScrapedAt:   time.Now().Unix(),
			WeightLimit: 300,
			Cooldown:    60,
		},
	}, nil
}

func (s *scraper) UpdateRepo(portfolio *model.Portfolio) error {
	repo, err := sqlite3.NewRepository(portfolio)
	if err != nil {
		return err
	}

	s.repo = repo

	if err := s.repo.UpdatePortfolios(config.Portfolios); err != nil {
		return err
	}

	return nil
}

func (s *scraper) UpdateExchange(portfolio *model.Portfolio) error {
	exchange, err := exchange.NewExchange(s.ctx.Portfolio, s.ctx)
	if err != nil {
		return err
	}

	s.exchange = exchange
	return nil
}

func (s *scraper) Scrape() error {
	for _, portfolio := range config.Portfolios {
		log.Printf("scraping portfolio %s", portfolio.Alias)
		if err := s.ScrapePortfolio(portfolio); err != nil {
			log.Print(fmt.Errorf("portfolio %s: %v", portfolio.Alias, err))
		}
	}

	return nil
}

func (s *scraper) ScrapePortfolio(portfolio *model.Portfolio) error {
	s.ctx.Portfolio = portfolio

	if err := s.UpdateRepo(portfolio); err != nil {
		return err
	}

	if err := s.UpdateExchange(portfolio); err != nil {
		return err
	}

	if err := s.repo.RemoveAllOrders(portfolio); err != nil {
		return err
	}

	tasks := []func() error{
		s.ScrapePositions,
		s.ScrapeOrders,
	}

	if portfolio.HistoryScraped {
		tasks = append(tasks, s.ScrapeIncome)
	} else {
		tasks = append(tasks, s.ScrapeIncomeHistory)
	}

	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}

		if s.IsWeightOverused() {
			log.Printf("weight limit exceeded, sleeping for %d seconds", s.ctx.Cooldown)
			s.Sleep()
		}
	}

	return nil
}

func (s *scraper) ScrapePositions() error {
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

func (s *scraper) ScrapeIncomeHistory() error {
	oldestIncomeTime := s.ctx.ScrapedAt

	for {
		if s.IsWeightOverused() {
			log.Printf("weight limit exceeded, sleeping for %d seconds", s.ctx.Cooldown)
			s.Sleep()
		}

		log.Printf("scraping income history from %d", oldestIncomeTime)
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

		oldestIncomeTime = incomes[0].Timestamp - 1
	}

	s.ctx.Portfolio.HistoryScraped = true
	if err := s.repo.UpdatePortfolio(s.ctx.Portfolio); err != nil {
		return err
	}

	return nil
}

func (s *scraper) IsWeightOverused() bool {
	return s.ctx.WeightUsed > s.ctx.WeightLimit
}

func (s *scraper) Sleep() {
	time.Sleep(time.Second * s.ctx.Cooldown)
	s.ctx.WeightUsed = 0
}
