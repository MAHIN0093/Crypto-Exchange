package orderbook

import (
	"fmt"
	"testing"
	"reflect"
)

func assert(t *testing.T, a, b any){
	if !reflect.DeepEqual(a, b){
		t.Errorf("Assertion failed: (%T: %+v) != (%T: %+v)", a, a, b, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10000)
	buyOrderA:= NewOrder(true, 5) 
	buyOrderB:= NewOrder(true, 8) 
	buyOrderC:= NewOrder(true, 1) 

	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)

	l.DeleteOrder(buyOrderB)

	fmt.Println(l) 
}
func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrderA:= NewOrder(false, 10)
	sellOrderB:= NewOrder(true, 5)
	ob.PlaceLimitOrder(10000, sellOrderA) 
	ob.PlaceLimitOrder(9000, sellOrderB)

	assert(t, len(ob.asks), 1)
	assert(t, len(ob.Orders), 2)
	assert(t, ob.Orders[sellOrderA.ID], sellOrderA)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderbook()

	sellOrder:= NewOrder(false, 20)
	ob.PlaceLimitOrder(10000, sellOrder)  

	buyOrder:= NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 10.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)

	assert(t, matches[0].SizeFilled, 10.0)
	assert(t, matches[0].Price, 10000.0)
	assert(t, buyOrder.IsFilled(), true)
	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultifill(t *testing.T) {
	ob := NewOrderbook()

	buyOrderA:= NewOrder(true, 5)
	buyOrderB:= NewOrder(true, 8)
	buyOrderC:= NewOrder(true, 10)
	buyOrderD:= NewOrder(true, 1)

	ob.PlaceLimitOrder(5000, buyOrderC)
	ob.PlaceLimitOrder(5000, buyOrderD)
	ob.PlaceLimitOrder(9000, buyOrderB)
	ob.PlaceLimitOrder(10000, buyOrderA)

	assert(t, ob.BidTotalVolume(), 24.0)

	sellOrder:= NewOrder(false, 20)
	matches := ob.PlaceMarketOrder(sellOrder)

	assert(t, ob.BidTotalVolume(), 4.0)
	assert(t, len(matches), 3) 
	assert(t, len(ob.bids), 1)

	fmt.Printf("%+v", matches)
}

func TestCancelorders(t *testing.T) {
	ob := NewOrderbook()

	buyOrder:= NewOrder(true, 5)
	ob.PlaceLimitOrder(10000, buyOrder)

	assert(t, ob.BidTotalVolume(), 5.0)

	ob.CancelOrder(buyOrder)

	assert(t, ob.BidTotalVolume(), 0.0)
	

	_, ok := ob.Orders[buyOrder.ID]
	assert(t, ok, false)
}