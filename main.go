package main

import (
	"fmt"
	"log"
	"scraper/scraper"

	"github.com/spf13/viper"
)

func main() {

	loadConfig()

	config := &scraper.ScraperConfig{
		APIKey:    viper.GetString("api_key"),
		APISecret: viper.GetString("api_secret"),
		DBPath:    viper.GetString("db_path"),
	}

	scraper, err := scraper.New(config)
	if err != nil {
		log.Fatal(err)
	}

	scraper.Scrape()

}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config: %v", err)
	}

	return nil
}
