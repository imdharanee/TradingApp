package engine

import (
	"sync"

	"dex/internal/models"
	"dex/internal/orderbook"
)

// MatchingEngine is responsible for matching orders and executing trades
type MatchingEngine struct {
	orderBook *orderbook.OrderBook
	orderChan chan *models.Order
	mutex     sync.Mutex
}

// NewMatchingEngine creates a new matching engine
func NewMatchingEngine(orderBook *orderbook.OrderBook) *MatchingEngine {
	return &MatchingEngine{
		orderBook: orderBook,
		orderChan: make(chan *models.Order, 100), // Buffer for 100 orders
	}
}

// Start begins processing orders
func (me *MatchingEngine) Start() {
	for order := range me.orderChan {
		me.processOrder(order)
	}
}

// Stop stops the matching engine
func (me *MatchingEngine) Stop() {
	close(me.orderChan)
}

// SubmitOrder adds an order to the processing queue
func (me *MatchingEngine) SubmitOrder(order *models.Order) {
	me.orderChan <- order
}

// processOrder handles the matching logic for a new order
func (me *MatchingEngine) processOrder(order *models.Order) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	// For market orders, we need to set the price based on the best available price
	if order.Type == models.MarketOrder {
		if order.Side == models.BuyOrder {
			sells := me.orderBook.GetSellOrders()
			if len(sells) > 0 {
				// Use the lowest sell price for market buy orders
				order.Price = sells[0].Price
			}
		} else {
			buys := me.orderBook.GetBuyOrders()
			if len(buys) > 0 {
				// Use the highest buy price for market sell orders
				order.Price = buys[0].Price
			}
		}
	}

	// Add the order to the book first
	me.orderBook.AddOrder(order)

	// Try to match the order
	if order.Side == models.BuyOrder {
		me.matchBuyOrder(order)
	} else {
		me.matchSellOrder(order)
	}

	// If the order is filled or cancelled, remove it from the book
	if order.Status == models.Filled || order.Status == models.Cancelled {
		me.orderBook.RemoveOrder(order.ID)
	}
}

// matchBuyOrder tries to match a buy order with existing sell orders
func (me *MatchingEngine) matchBuyOrder(buyOrder *models.Order) {
	sellOrders := me.orderBook.GetSellOrders()

	for _, sellOrder := range sellOrders {
		// Stop if the buy order is filled
		if buyOrder.Remaining <= 0 {
			break
		}

		// Skip if the sell price is higher than the buy price
		if sellOrder.Price > buyOrder.Price {
			continue
		}

		// Calculate the quantity to trade
		quantity := min(buyOrder.Remaining, sellOrder.Remaining)

		// Execute the trade
		buyOrder.Fill(quantity)
		sellOrder.Fill(quantity)

		// Create a trade record
		trade := models.NewTrade(buyOrder.ID, sellOrder.ID, sellOrder.Price, quantity)
		me.orderBook.AddTrade(trade)

		// Remove the sell order if it's filled
		if sellOrder.Status == models.Filled {
			me.orderBook.RemoveOrder(sellOrder.ID)
		}
	}
}

// matchSellOrder tries to match a sell order with existing buy orders
func (me *MatchingEngine) matchSellOrder(sellOrder *models.Order) {
	buyOrders := me.orderBook.GetBuyOrders()

	for _, buyOrder := range buyOrders {
		// Stop if the sell order is filled
		if sellOrder.Remaining <= 0 {
			break
		}

		// Skip if the buy price is lower than the sell price
		if buyOrder.Price < sellOrder.Price {
			continue
		}

		// Calculate the quantity to trade
		quantity := min(sellOrder.Remaining, buyOrder.Remaining)

		// Execute the trade
		sellOrder.Fill(quantity)
		buyOrder.Fill(quantity)

		// Create a trade record
		trade := models.NewTrade(buyOrder.ID, sellOrder.ID, buyOrder.Price, quantity)
		me.orderBook.AddTrade(trade)

		// Remove the buy order if it's filled
		if buyOrder.Status == models.Filled {
			me.orderBook.RemoveOrder(buyOrder.ID)
		}
	}
}

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
