package exchange

import (
	"context"
	"fmt"
	"log"
	"scraper/scraper/models"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
)

type binanceFutures struct {
	portfolio *models.Portfolio
	client    *futures.Client
}

func NewBinanceFutures(portfolio *models.Portfolio) (Exchange, error) {

	client := futures.NewClient(portfolio.APIKey, portfolio.APISecret)
	err := client.NewPingService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping binance api: %v", err)
	}

	return &binanceFutures{
		portfolio: portfolio,
		client:    client,
	}, nil
}

func (e *binanceFutures) GetPositions() ([]*models.Position, error) {
	account, err := e.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var positionsModels []*models.Position
	for _, position := range account.Positions {

		positionCost, err := e.getPositionCost(position)
		if err != nil {
			return nil, err
		}

		if positionCost < 0.1 {
			continue
		}

		model, err := e.parsePosition(position)
		if err != nil {
			return nil, err
		}

		positionsModels = append(positionsModels, model)
	}
	return positionsModels, nil
}

func (e *binanceFutures) GetOrders() ([]*models.Order, error) {
	orders, err := e.client.NewListOpenOrdersService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	var ordersModels []*models.Order
	for _, order := range orders {
		model, err := e.parseOrder(order)
		if err != nil {
			return nil, err
		}

		ordersModels = append(ordersModels, model)
	}

	return ordersModels, nil
}

func (e *binanceFutures) GetTrades(symbol string) ([]*models.Trade, error) {
	// return e.client.NewRecentTradesService().Symbol(symbol).Do(context.Background())
	return nil, nil
}

func (e *binanceFutures) getPositionCost(position *futures.AccountPosition) (float64, error) {
	cost, err := strconv.ParseFloat(position.PositionInitialMargin, 64)
	if err != nil {
		return 0, err
	}

	return cost, nil
}

func (e *binanceFutures) parsePosition(ap *futures.AccountPosition) (*models.Position, error) {
	amount, err := strconv.ParseFloat(ap.PositionAmt, 64)
	if err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(ap.PositionInitialMargin, 64)
	if err != nil {
		return nil, err
	}

	unpnl, err := strconv.ParseFloat(ap.UnrealizedProfit, 64)
	if err != nil {
		return nil, err
	}

	lev, err := strconv.ParseInt(ap.Leverage, 10, 32)
	if err != nil {
		return nil, err
	}

	return &models.Position{
		Symbol:    ap.Symbol,
		Amount:    amount,
		Cost:      price,
		Isolated:  ap.Isolated,
		UnPnl:     unpnl,
		Side:      string(ap.PositionSide),
		Leverage:  int32(lev),
		UpdatedAt: ap.UpdateTime,
	}, nil
}

func (e *binanceFutures) parseOrder(order *futures.Order) (*models.Order, error) {

	price, err := strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(order.OrigQuantity, 64)
	if err != nil {
		return nil, err
	}

	return &models.Order{
		Symbol: order.Symbol,
		Price:  price,
		Amount: amount,
	}, nil
}

func (e *binanceFutures) parseTrade(trade *futures.Trade) (*models.Trade, error) {
	// return &models.Trade{
	// 	Symbol:    trade.Symbol,
	// 	Price:     trade.Price,
	// 	Amount:    trade.Quantity,
	// 	Timestamp: trade.Time,
	// }, nil

	return nil, nil
}
