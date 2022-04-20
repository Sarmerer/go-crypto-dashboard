package main

import (
	"fmt"
	"log"

	"github.com/sarmerer/go-crypto-dashboard/config"
	"github.com/sarmerer/go-crypto-dashboard/tracker/repository"
	"github.com/sarmerer/go-crypto-dashboard/tracker/repository/sqlite3"
	"github.com/sarmerer/go-crypto-dashboard/tracker/scraper"
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

	StartScraper()

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
