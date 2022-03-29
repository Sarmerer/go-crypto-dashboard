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
	SERVE
)

func main() {

	err := config.Load()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to load config: %v", err))
	}

	args := os.Args[1:]
	command := SCRAPE
	if len(args) > 0 {
		if args[0] == "serve" {
			command = SERVE
		}
	}

	switch command {
	case SCRAPE:
		scraper, err := scraper.NewScraper()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to initialize scraper: %v", err))
		}

		scraper.ContinuousScrape()
	case SERVE:
		api.Serve()
		return
	default:
		log.Fatal(fmt.Errorf("unknown command: %v", command))
	}

}
