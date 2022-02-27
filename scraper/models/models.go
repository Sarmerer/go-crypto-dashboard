package models

type Portfolio struct {
	ID        int32  `gorm:"primaryKey" mapstructure:"id"`
	Name      string `gorm:"type:varchar(50);not null" mapstructure:"name"`
	Exchange  string `gorm:"type:varchar(50);not null" mapstructure:"exchange"`
	APIKey    string `gorm:"-" mapstructure:"api_key"`
	APISecret string `gorm:"-" mapstructure:"api_secret"`

	DBPath string `gorm:"-" mapstructure:"db_path"`
}

type ScrapeCtx struct {
	ScrapedAt   int64      `gorm:"type:bigint;not null"`
	Portfolio   *Portfolio `gorm:"foreignKey:PortfolioID;reference:ID"`
	PortfolioID int32      `gorm:"type:int;not null"`
}

type Order struct {
	ScrapeCtx

	Symbol string
	Price  float64
	Amount float64
}

type Trade struct {
	ScrapeCtx

	Symbol string
	Price  float64
	Amount float64
}

type SymbolPrice struct {
	Symbol string
	Price  float64
}
