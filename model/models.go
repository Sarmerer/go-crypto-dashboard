package model

import (
	"time"
)

type PortfolioID string

type Portfolio struct {
	ID             PortfolioID `gorm:"primaryKey;varchar(50);not null" mapstructure:"id"`
	Alias          string      `gorm:"type:varchar(50);not null" mapstructure:"alias"`
	Exchange       string      `gorm:"type:varchar(50);not null" mapstructure:"exchange"`
	APIKey         string      `gorm:"-" mapstructure:"key"`
	APISecret      string      `gorm:"-" mapstructure:"secret"`
	HistoryScraped bool        `gorm:"type:bool;default:false" mapstructure:"-"`
}

func (p *Portfolio) SyncWith(record *Portfolio) {
	p.HistoryScraped = record.HistoryScraped
}

type ScrapeCtx struct {
	ScrapedAt   int64       `gorm:"type:bigint;not null"`
	Portfolio   *Portfolio  `gorm:"foreignkey:PortfolioID"`
	PortfolioID PortfolioID `gorm:"type:varchar(50);not null"`

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

	ID       uint      `gorm:"primaryKey"`
	Symbol   string    `gorm:"type:varchar(20);not null"`
	Side     string    `gorm:"type:varchar(7);not null"`
	Amount   float64   `gorm:"type:float;not null"`
	Cost     float64   `gorm:"type:float;not null"`
	Isolated bool      `gorm:"type:bool;not null"`
	UnPnl    float64   `gorm:"type:float;not null"`
	Leverage int32     `gorm:"type:int;not null"`
	Date     time.Time `gorm:"type:date;not null"`
}

type Order struct {
	ScrapeCtx

	ID           int64     `gorm:"primaryKey; autoIncrement:false; type:bigint; not null"`
	Symbol       string    `gorm:"type:varchar(20);not null"`
	Side         string    `gorm:"type:varchar(7);not null"`
	PositionSide string    `gorm:"type:varchar(7);not null"`
	TimeInForce  string    `gorm:"type:varchar(7);not null"`
	Type         string    `gorm:"type:varchar(7);not null"`
	Price        float64   `gorm:"type:float;not null"`
	Amount       float64   `gorm:"type:float;not null"`
	ReduceOnly   bool      `gorm:"type:bool;not null"`
	Date         time.Time `gorm:"type:date;not null"`
}

type Income struct {
	ScrapeCtx

	ID      int64     `gorm:"primaryKey; autoIncrement:false; type:bigint; not null"`
	Symbol  string    `gorm:"type:varchar(20); not null"`
	Asset   string    `gorm:"type:varchar(20); not null"`
	Type    string    `gorm:"type:varchar(20); not null"`
	Info    string    `gorm:"type:varchar(20); not null"`
	Income  float64   `gorm:"type:float; not null"`
	TradeID int64     `gorm:"type:bigint; not null"`
	Date    time.Time `gorm:"type:date; not null"`
}

type DailyBalance struct {
	ScrapeCtx

	ID      uint      `gorm:"primaryKey"`
	Balance float64   `gorm:"type:float; not null"`
	Date    time.Time `gorm:"type:date; not null"`
}
