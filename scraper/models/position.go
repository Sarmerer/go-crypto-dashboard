package models

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

func (p *Position) ApplyScrapeCtx(ctx *ScrapeCtx) {
	p.ScrapeCtx = *ctx
}
