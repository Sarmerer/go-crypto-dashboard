package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sarmerer/go-crypto-dashboard/api"
	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/scraper"
)

const (
	SCRAPE = iota
	INFO
)

func main() {

	err := config.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %v", err))
	}

	args := os.Args[1:]
	command := SCRAPE
	if len(args) > 0 {
		if args[0] == "info" {
			command = INFO
		}
	}

	switch command {
	case SCRAPE:
		scraper, err := scraper.NewScraper()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to initialize scraper: %v", err))
		}

		if err := scraper.Scrape(); err != nil {
			log.Fatal(fmt.Errorf("failed to scrape: %v", err))
		}
	case INFO:
		api.PortfolioInfo(config.Portfolios[0])
		return
	default:
		log.Fatal(fmt.Errorf("unknown command: %v", command))
	}

}
