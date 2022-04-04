package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sarmerer/go-crypto-dashboard/api"
	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/repository"
	"github.com/sarmerer/go-crypto-dashboard/repository/sqlite3"
	"github.com/sarmerer/go-crypto-dashboard/scraper"
	"gorm.io/driver/sqlite"
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
	command := GetCommand(args)

	switch command {
	case SCRAPE:
		StartScraper()
	case SERVE:
		StartAPI()
	default:
		log.Fatal(fmt.Errorf("unknown command: %s", args[0]))
	}
}

func GetCommand(args []string) (command int) {
	if len(args) == 0 {
		return SCRAPE
	}

	switch args[0] {
	case "serve":
		return SERVE
	default:
		return -1
	}
}

func GetRepo() (repo repository.Repository, err error) {
	dialector := sqlite.Open(config.DBPath)
	repo, err = sqlite3.NewRepository(dialector)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %v", err)
	}

	return repo, nil
}

func StartScraper() {
	repo, err := GetRepo()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize repository: %v", err))
	}

	scraper, err := scraper.NewScraper(repo)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize scraper: %v", err))
	}

	err = scraper.ContinuousScrape()
	if err != nil {
		log.Fatal(fmt.Errorf("initial scrape failed: %v", err))
	}
}

func StartAPI() {
	repo, err := GetRepo()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize repository: %v", err))
	}

	api.Serve(repo)
}
