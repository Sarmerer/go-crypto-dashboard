package api

import (
	"fmt"
	"time"

	"github.com/sarmerer/go-crypto-dashboard/database/sqlite3"
	"github.com/sarmerer/go-crypto-dashboard/model"
)

// function that prints following info:
// - open positions
// - amount of open buy and sell orders for each symbol
// - total income for the past day and month
func PortfolioInfo(portfolio *model.Portfolio) {
	fmt.Printf("Portfolio: %s\n", portfolio.Alias)

	repo, err := sqlite3.NewRepository(portfolio)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	positions, err := repo.GetPositions()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Open positions:\n")
	for _, position := range positions {
		fmt.Printf("%s: %f\n", position.Symbol, position.Amount)
	}

	orders, err := repo.GetOrders()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("Open orders:\n")
	for _, order := range orders {
		fmt.Printf("%s %s: %f\n", order.Side, order.Symbol, order.Amount)
	}

	start := time.Now().AddDate(0, 0, -1).UnixMilli()
	end := time.Now().UnixMilli()
	income, err := repo.GetIncomeBetween(start, end)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	sum := float64(0)
	for _, inc := range income {
		sum += inc.Income
	}
	fmt.Printf("Yesterday's income: %f\n", sum)

	start = time.Now().AddDate(0, -1, 0).UnixMilli()
	end = time.Now().UnixMilli()
	income, err = repo.GetIncomeBetween(start, end)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	sum = float64(0)
	for _, inc := range income {
		sum += inc.Income
	}

	fmt.Printf("Last month's income: %f\n", sum)

}
