package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/MAHIN0093/go-lang/orderbook"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	ex := NewExchange()

	e.GET("/book/:market", ex.handleGetbook)

	e.POST("/order", ex.handlePlaceOrder)

	e.DELETE("/order/:id", ex.CancelOrder)

	e.Start(":3000")

}

type Market string

const (
	MarketETH Market = "ETH"
)

type Ordertype string

const (
	MarketOrder Ordertype = "MARKET"
	LimitOrder  Ordertype = "LIMIT"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {

	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type PlaceOrderRequest struct {
	Type   Ordertype
	Bid    bool
	Size   float64
	Price  float64
	Market Market
}

type Order struct {
	ID        int64
	Price     float64
	Size      float64
	Bid       bool
	Timestamp int64
}

type orderbookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*Order
	Bids           []*Order
}

func (ex *Exchange) handleGetbook(c echo.Context) error {
	market := Market(c.Param("market"))

	ob, ok := ex.orderbooks[market]

	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"message": "Market not found",
		})
	}

	orderbookData := orderbookData{
		Asks:           []*Order{},
		Bids:           []*Order{},
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Bids = append(orderbookData.Bids, &o)
		}
	}

	return c.JSON(http.StatusOK, orderbookData)
}

func (ex *Exchange) CancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	ob := ex.orderbooks[MarketETH]
	orderCanceled := false

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			if order.ID == int64(id) {
				ob.CancelOrder(order)
				orderCanceled = true
			}

			if orderCanceled {
				return c.JSON(http.StatusOK, map[string]any{"message": "Order canceled successfully"})
			}
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			if order.ID == int64(id) {
				ob.CancelOrder(order)
				orderCanceled = true
			}

			if orderCanceled {
				return c.JSON(http.StatusOK, map[string]any{"message": "Order canceled successfully"})
			}
		}
	}

	return nil
}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest

	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)

	ob := ex.orderbooks[market]

	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size)

	if placeOrderData.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderData.Price, order)
		return c.JSON(200, map[string]any{"message": "Limit Order placed successfully"})
	}

	if placeOrderData.Type == MarketOrder {
		matches := ob.PlaceMarketOrder(order)
		return c.JSON(200, map[string]any{"matches": len(matches)})
	}

	return nil
}
