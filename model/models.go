package model

import (
	"time"
)

type Portfolio struct {
	ID             int64  `gorm:"primary_key;varchar(50);not null" mapstructure:"id"`
	Alias          string `gorm:"type:varchar(50);not null" mapstructure:"alias"`
	Exchange       string `gorm:"type:varchar(50);not null" mapstructure:"exchange"`
	APIKey         string `gorm:"-" mapstructure:"api_key"`
	APISecret      string `gorm:"-" mapstructure:"api_secret"`
	HistoryScraped bool   `gorm:"type:boolean;default:false" mapstructure:"-"`
}

type ScrapeCtx struct {
	ScrapedAt   int64      `gorm:"type:bigint;not null"`
	Portfolio   *Portfolio `gorm:"foreignKey:PortfolioID;reference:ID"`
	PortfolioID int64      `gorm:"type:int;not null"`

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

	Symbol    string  `gorm:"primaryKey;type:varchar(20);not null"`
	Side      string  `gorm:"primaryKey;type:varchar(7);not null"`
	Amount    float64 `gorm:"type:float;not null"`
	Cost      float64 `gorm:"type:float;not null"`
	Isolated  bool    `gorm:"type:bool;not null"`
	UnPnl     float64 `gorm:"type:float;not null"`
	Leverage  int32   `gorm:"type:int;not null"`
	UpdatedAt int64   `gorm:"type:int;not null"`
}

type Order struct {
	ScrapeCtx

	ID           int64   `gorm:"primaryKey;type:varchar(20);not null"`
	Symbol       string  `gorm:"type:varchar(20);not null"`
	Side         string  `gorm:"type:varchar(7);not null"`
	PositionSide string  `gorm:"type:varchar(7);not null"`
	TimeInForce  string  `gorm:"type:varchar(7);not null"`
	Type         string  `gorm:"type:varchar(7);not null"`
	Price        float64 `gorm:"type:float;not null"`
	Amount       float64 `gorm:"type:float;not null"`
	ReduceOnly   bool    `gorm:"type:bool;not null"`
	Timestamp    int64   `gorm:"type:int;not null"`
}

type Income struct {
	ScrapeCtx

	ID        int64   `gorm:"primaryKey; type:bigint"`
	Symbol    string  `gorm:"type:varchar(20);not null"`
	Asset     string  `gorm:"type:varchar(20);not null"`
	Type      string  `gorm:"type:varchar(20);not null"`
	Info      string  `gorm:"type:varchar(20);not null"`
	Income    float64 `gorm:"type:float;not null"`
	Timestamp int64   `gorm:"type:bigint;not null"`
	TradeID   int64   `gorm:"type:bigint;not null"`
}

type IncomeSummary struct {
	Symbol    string  `gorm:"primaryKey;type:varchar(20);not null"`
	Income    float64 `gorm:"type:float;not null"`
	Timestamp int64   `gorm:"type:bigint;not null"`
}
