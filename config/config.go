package config

import (
	"fmt"
	"time"

	"github.com/sarmerer/go-crypto-dashboard/tracker/model"
	"github.com/spf13/viper"
)

const (
	DefaultAPIPort int    = 3000
	DefaultDBPath  string = "./database.db"

	DefaultScrapeHistory  bool          = true
	DefaultScrapeInterval time.Duration = time.Minute * 5

	DefaultExcWeightLimit    int32         = 500
	DefaultExcWeightCooldown time.Duration = 60 * time.Second
)

var (
	APIPort int    = DefaultAPIPort
	DBPath  string = DefaultDBPath

	ScrapeHistory                = DefaultScrapeHistory
	ScrapeInterval time.Duration = DefaultScrapeInterval

	ExchangeWeightLimit    int32         = DefaultExcWeightLimit
	ExchangeWeightCooldown time.Duration = DefaultExcWeightCooldown
)

func Load() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	var cooldown int64
	var interval int64
	fields := map[string]interface{}{
		"api_port": &APIPort,

		"scrape_history":       &ScrapeHistory,
		"scrape_interval_secs": &interval,

		"exchange_weight_limit":  &ExchangeWeightLimit,
		"exchange_cooldown_secs": &cooldown,
	}

	for field, ptr := range fields {
		if err := viper.UnmarshalKey(field, ptr); err != nil {
			return fmt.Errorf("failed to unmarshal key %s: %v", field, err)
		}
	}

	ExchangeWeightCooldown = time.Duration(cooldown) * time.Second
	ScrapeInterval = time.Duration(interval) * time.Second

	return nil
}

func GetPortfolios() ([]*model.Portfolio, error) {
	portfolios := []*model.Portfolio{}

	if err := viper.UnmarshalKey("portfolios", &portfolios); err != nil {
		return nil, fmt.Errorf("failed to unmarshal portfolios: %v", err)
	}

	return portfolios, nil
}
