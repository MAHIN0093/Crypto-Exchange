package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"crypto/ecdsa"
	"math/big"

	"github.com/MAHIN0093/go-lang/orderbook"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

const (
	exchangePrivateKey           = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	userPrivateKey               = "829e924fdf021ba3dbbc4225edfece9aca04b929d6e75613329ca6f1d31c0bb4"
	MarketOrder        Ordertype = "MARKET"
	LimitOrder         Ordertype = "LIMIT"
	MarketETH          Market    = "ETH"
)

type (
	Market    string
	Ordertype string

	MatchedOrder struct {
		UserID int64
		Size  float64
		Price float64
		ID    int64
	}

	PlaceOrderRequest struct {
		UserID int64
		Type   Ordertype
		Bid    bool
		Size   float64
		Price  float64
		Market Market
	}

	Order struct {
		UserID    int64
		ID        int64
		Price     float64
		Size      float64
		Bid       bool
		Timestamp int64
	}

	orderbookData struct {
		TotalBidVolume float64
		TotalAskVolume float64
		Asks           []*Order
		Bids           []*Order
	}
)

func main() {
	e := echo.New()

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}

	e.HTTPErrorHandler = httpErrorHandler

	ex, err := NewExchange(exchangePrivateKey, client)
	if err != nil {
		log.Fatal(err)
	}

	exAcc := NewUser(0, exchangePrivateKey)
	userAcc := NewUser(8, userPrivateKey)

	ex.Users[exAcc.ID] = exAcc
	ex.Users[userAcc.ID] = userAcc

	e.GET("/book/:market", ex.handleGetbook)

	e.POST("/order", ex.handlePlaceOrder)

	e.DELETE("/order/:id", ex.CancelOrder)

	exAddress := crypto.PubkeyToAddress(exAcc.PrivateKey.PublicKey)
	exBalance, _ := client.BalanceAt(context.Background(), exAddress, nil)
	fmt.Printf("Exchange Balance: %d\n", exBalance)

	userAddress := crypto.PubkeyToAddress(userAcc.PrivateKey.PublicKey)
	userBalance, _ := client.BalanceAt(context.Background(), userAddress, nil)
	fmt.Printf("User Balance: %d\n", userBalance)

	e.Start(":3000")

}

type User struct {
	ID         int64
	PrivateKey *ecdsa.PrivateKey
}

func NewUser(id int64, privateKey string) *User {
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	return &User{
		ID:         id,
		PrivateKey: pk,
	}
}

func httpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}

type Exchange struct {
	Client     *ethclient.Client
	Users      map[int64]*User
	orders     map[int64]int64
	PrivateKey *ecdsa.PrivateKey
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange(privateKey string, client *ethclient.Client) (*Exchange, error) {

	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	pk, err := crypto.HexToECDSA(exchangePrivateKey)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		Client:     client,
		Users:      make(map[int64]*User),
		orders:     make(map[int64]int64),
		PrivateKey: pk,
		orderbooks: orderbooks,
	}, nil
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
				UserID:    order.UserID,
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
				UserID:    order.UserID,
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
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	return c.JSON(200, map[string]any{"message": "Limit Order Deleted successfully"})
}

func (ex *Exchange) handlePlaceMarketOrder(market Market, order *orderbook.Order) ([]orderbook.Match, []*MatchedOrder) {
	ob := ex.orderbooks[market]
	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := false
	if order.Bid {
		isBid = true
	}

	for i := range matchedOrders {
		id := matches[i].Bid.ID
		userID := matches[i].Bid.UserID
		if isBid {
			id = matches[i].Ask.ID
			userID = matches[i].Ask.UserID
		}
		matchedOrders[i] = &MatchedOrder{
			UserID: userID,
			ID:    id,
			Size:  matches[i].SizeFilled,
			Price: matches[i].Price,
		}
	}

	return matches, matchedOrders

}

func (ex *Exchange) handlePlaceLimitOrder(market Market, price float64, order *orderbook.Order) error {
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order)

	return nil

}

func (ex *Exchange) handlePlaceOrder(c echo.Context) error {

	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size, placeOrderData.UserID)

	if placeOrderData.Type == LimitOrder {
		if err := ex.handlePlaceLimitOrder(market, placeOrderData.Price, order); err != nil {
			fmt.Println(err)
		}
		return c.JSON(200, map[string]any{"message": "Limit Order placed successfully"})
	}

	if placeOrderData.Type == MarketOrder {
		matches, matchedOrders := ex.handlePlaceMarketOrder(market, order)

		if err := ex.handleMatches(matches); err != nil {
			return err
		}

		return c.JSON(200, map[string]any{"matches": matchedOrders})
	}

	return nil
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	for _, match := range matches {
		fromUser, ok := ex.Users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Ask.UserID)
		}

		toUser, ok := ex.Users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Bid.UserID)
		}
		toAddress := crypto.PubkeyToAddress(toUser.PrivateKey.PublicKey)

		amount := big.NewInt(int64(match.SizeFilled))

		transferETH(ex.Client, fromUser.PrivateKey, toAddress, amount)

	}
	return nil
}
