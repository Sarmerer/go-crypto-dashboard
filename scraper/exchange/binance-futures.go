package exchange

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/sarmerer/go-crypto-dashboard/model"
)

type binanceFutures struct {
	portfolio *model.Portfolio
	client    *futures.Client
	ctx       *model.ScrapeCtx

	UnderlyingTransport http.RoundTripper
}

// RoundTrip implement http roundtrip
func (e *binanceFutures) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := e.UnderlyingTransport.RoundTrip(req)
	if resp != nil && resp.Header != nil {
		weight, err := strconv.Atoi(resp.Header.Get("X-Mbx-Used-Weight-1m"))
		if err != nil {
			return resp, nil
		}

		e.ctx.WeightUsed = int32(weight)
	}

	return resp, err
}

func NewBinanceFutures(portfolio *model.Portfolio, ctx *model.ScrapeCtx) (Exchange, error) {
	client := futures.NewClient(portfolio.APIKey, portfolio.APISecret)
	err := client.NewPingService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping binance api: %v", err)
	}

	client.HTTPClient = &http.Client{Transport: &binanceFutures{
		UnderlyingTransport: http.DefaultTransport,
		ctx:                 ctx,
	}}

	exchange := &binanceFutures{
		portfolio: portfolio,
		client:    client,
		ctx:       ctx,

		UnderlyingTransport: http.DefaultTransport,
	}

	return exchange, nil
}

func (e *binanceFutures) GetPositions() ([]*model.Position, error) {
	account, err := e.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %v", err)
	}

	var positions []*model.Position
	for _, rawPosition := range account.Positions {

		positionCost, err := e.getPositionCost(rawPosition)
		if err != nil {
			return nil, err
		}

		if positionCost < 0.1 {
			continue
		}

		position, err := e.parsePosition(rawPosition)
		if err != nil {
			return nil, err
		}

		positions = append(positions, position)
	}
	return positions, nil
}

func (e *binanceFutures) GetOrders() ([]*model.Order, error) {
	rawOrders, err := e.client.NewListOpenOrdersService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	for _, rawOrder := range rawOrders {
		order, err := e.parseOrder(rawOrder)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}
	return orders, nil
}

func (e *binanceFutures) GetIncome() ([]*model.Income, error) {
	service := e.client.NewGetIncomeHistoryService()

	rawIncomes, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	var incomes []*model.Income
	for _, rawIncome := range rawIncomes {
		income, err := e.parseIncome(rawIncome)
		if err != nil {
			return nil, err
		}

		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (e *binanceFutures) GetIncomeBetween(startTime, endTime int64) ([]*model.Income, error) {
	service := e.client.NewGetIncomeHistoryService()

	if startTime > 0 {
		service.StartTime(startTime)
	}

	if endTime > 0 {
		service.EndTime(endTime)
	}

	rawIncomes, err := service.Do(context.Background())
	if err != nil {
		return nil, err
	}

	var incomes []*model.Income
	for _, rawIncome := range rawIncomes {
		income, err := e.parseIncome(rawIncome)
		if err != nil {
			return nil, err
		}

		incomes = append(incomes, income)
	}

	return incomes, nil
}

func (e *binanceFutures) getPositionCost(position *futures.AccountPosition) (float64, error) {
	cost, err := strconv.ParseFloat(position.PositionInitialMargin, 64)
	if err != nil {
		return 0, err
	}

	return cost, nil
}

func (e *binanceFutures) parsePosition(ap *futures.AccountPosition) (*model.Position, error) {
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

	return &model.Position{
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

func (e *binanceFutures) parseOrder(order *futures.Order) (*model.Order, error) {
	price, err := strconv.ParseFloat(order.Price, 64)
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(order.OrigQuantity, 64)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:           order.OrderID,
		Symbol:       order.Symbol,
		Side:         string(order.Side),
		PositionSide: string(order.PositionSide),
		TimeInForce:  string(order.TimeInForce),
		Type:         string(order.Type),
		Price:        price,
		Amount:       amount,
		ReduceOnly:   order.ReduceOnly,
		Timestamp:    order.Time,
	}, nil
}

func (e *binanceFutures) parseIncome(income *futures.IncomeHistory) (*model.Income, error) {
	pnl, err := strconv.ParseFloat(income.Income, 64)
	if err != nil {
		return nil, err
	}

	var tradeID int64 = -1
	if income.TradeID != "" {
		tradeID, err = strconv.ParseInt(income.TradeID, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &model.Income{
		ID:        income.TranID,
		Symbol:    income.Symbol,
		Asset:     income.Asset,
		Type:      string(income.IncomeType),
		Info:      income.Info,
		Income:    pnl,
		Timestamp: income.Time,
		TradeID:   tradeID,
	}, nil
}
