package main

import (
	"fmt"
	"log"
	"scraper/scraper"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(fmt.Errorf("failed to read config: %v", err))
	}

	config := &scraper.ScraperConfig{
		DBPath: viper.GetString("db_path"),
	}

	if err := viper.UnmarshalKey("portfolios", &config.Portfolios); err != nil {
		log.Fatal(err)
	}

	scraper, err := scraper.NewScraper(config)
	if err != nil {
		log.Fatal(err)
	}

	scraper.Scrape()
}
