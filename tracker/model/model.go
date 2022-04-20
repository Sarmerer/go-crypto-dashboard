package model

import (
	"math"
	"time"
)

type PortfolioID string

type Portfolio struct {
	ID             PortfolioID `gorm:"primaryKey;varchar(50)" mapstructure:"id"`
	Alias          string      `gorm:"type:varchar(50)" mapstructure:"alias"`
	Exchange       string      `gorm:"type:varchar(50)" mapstructure:"exchange"`
	APIKey         string      `gorm:"-" mapstructure:"key"`
	APISecret      string      `gorm:"-" mapstructure:"secret"`
	HistoryScraped bool        `gorm:"type:bool;default:false" mapstructure:"-"`
}

func (p *Portfolio) SyncWith(record *Portfolio) {
	p.HistoryScraped = record.HistoryScraped
}

type ScrapeCtx struct {
	ScrapedAt   int64       `gorm:"type:bigint"`
	Portfolio   *Portfolio  `gorm:"foreignkey:PortfolioID"`
	PortfolioID PortfolioID `gorm:"type:varchar(50)"`

	WeightUsed  int32         `gorm:"-"`
	WeightLimit int32         `gorm:"-"`
	Cooldown    time.Duration `gorm:"-"`
}

func (p *ScrapeCtx) Apply(ctx *ScrapeCtx) {
	p.ScrapedAt = ctx.ScrapedAt
	p.Portfolio = ctx.Portfolio
	p.PortfolioID = ctx.PortfolioID

	p.WeightUsed = ctx.WeightUsed
}

type Position struct {
	ScrapeCtx

	ID         uint      `gorm:"primaryKey"`
	Symbol     string    `gorm:"type:varchar(20)"`
	Side       string    `gorm:"type:varchar(7)"`
	Amount     float64   `gorm:"type:float"`
	Cost       float64   `gorm:"type:float"`
	EntryPrice float64   `gorm:"type:float"`
	Isolated   bool      `gorm:"type:bool"`
	UnPnl      float64   `gorm:"type:float"`
	Leverage   int32     `gorm:"type:int"`
	Date       time.Time `gorm:"type:date"`
}

func (p *Position) IsOpen() bool {
	return math.Abs(p.Amount) > 0.0
}

type Order struct {
	ScrapeCtx

	ID           int64     `gorm:"primaryKey; autoIncrement:false; type:bigint"`
	Symbol       string    `gorm:"type:varchar(20)"`
	Side         string    `gorm:"type:varchar(7)"`
	PositionSide string    `gorm:"type:varchar(7)"`
	TimeInForce  string    `gorm:"type:varchar(7)"`
	Type         string    `gorm:"type:varchar(7)"`
	Price        float64   `gorm:"type:float"`
	Amount       float64   `gorm:"type:float"`
	ReduceOnly   bool      `gorm:"type:bool"`
	Date         time.Time `gorm:"type:date"`
}

type Income struct {
	ScrapeCtx

	ID      int64     `gorm:"primaryKey; autoIncrement:false; type:bigint"`
	Type    string    `gorm:"primaryKey; type:varchar(20)"`
	Symbol  string    `gorm:"type:varchar(20)"`
	Asset   string    `gorm:"type:varchar(20)"`
	Info    string    `gorm:"type:varchar(20)"`
	Income  float64   `gorm:"type:float"`
	TradeID int64     `gorm:"type:bigint"`
	Date    time.Time `gorm:"type:date"`
}

type DailyBalance struct {
	ScrapeCtx

	ID      uint      `gorm:"primaryKey"`
	Balance float64   `gorm:"type:float"`
	Date    time.Time `gorm:"type:date"`
}

type CurrentBalance struct {
	ScrapeCtx

	ID      uint      `gorm:"primaryKey"`
	Balance float64   `gorm:"type:float"`
	Date    time.Time `gorm:"type:date"`
}

type SymbolPrice struct {
	ScrapeCtx

	Symbol string  `gorm:"primaryKey;type:varchar(20)"`
	Price  float64 `gorm:"type:float"`
}
