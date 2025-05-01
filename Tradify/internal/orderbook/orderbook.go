package orderbook

import (
	"sort"
	"sync"

	"dex/internal/models"
)

type OrderBook struct {
	buys      []*models.Order          // buy orders sorted by price (highest first)
	sells     []*models.Order          // sell orders sorted by price (lowest first)
	ordersMap map[string]*models.Order // quick lookup by order ID
	trades    []*models.Trade          // record of executed trades
	mutex     sync.RWMutex
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		buys:      make([]*models.Order, 0),
		sells:     make([]*models.Order, 0),
		ordersMap: make(map[string]*models.Order),
		trades:    make([]*models.Trade, 0),
	}
}

func (ob *OrderBook) AddOrder(order *models.Order) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	// Add to map for quick lookup
	ob.ordersMap[order.ID] = order

	// Add to appropriate side and sort
	if order.Side == models.BuyOrder {
		ob.buys = append(ob.buys, order)
		// Sort buy orders by price (highest first)
		sort.Slice(ob.buys, func(i, j int) bool {
			return ob.buys[i].Price > ob.buys[j].Price
		})
	} else {
		ob.sells = append(ob.sells, order)
		// Sort sell orders by price (lowest first)
		sort.Slice(ob.sells, func(i, j int) bool {
			return ob.sells[i].Price < ob.sells[j].Price
		})
	}
}

func (ob *OrderBook) RemoveOrder(orderID string) bool {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	// Check if order exists
	order, exists := ob.ordersMap[orderID]
	if !exists {
		return false
	}

	// Remove from map
	delete(ob.ordersMap, orderID)

	// Remove from appropriate side
	if order.Side == models.BuyOrder {
		for i, o := range ob.buys {
			if o.ID == orderID {
				ob.buys = append(ob.buys[:i], ob.buys[i+1:]...)
				break
			}
		}
	} else {
		for i, o := range ob.sells {
			if o.ID == orderID {
				ob.sells = append(ob.sells[:i], ob.sells[i+1:]...)
				break
			}
		}
	}

	return true
}

func (ob *OrderBook) GetOrder(orderID string) (*models.Order, bool) {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	order, exists := ob.ordersMap[orderID]
	return order, exists
}

func (ob *OrderBook) GetBuyOrders() []*models.Order {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	result := make([]*models.Order, len(ob.buys))
	copy(result, ob.buys)
	return result
}

func (ob *OrderBook) GetSellOrders() []*models.Order {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	result := make([]*models.Order, len(ob.sells))
	copy(result, ob.sells)
	return result
}

func (ob *OrderBook) AddTrade(trade *models.Trade) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	ob.trades = append(ob.trades, trade)
}

func (ob *OrderBook) GetTrades() []*models.Trade {
	ob.mutex.RLock()
	defer ob.mutex.RUnlock()

	result := make([]*models.Trade, len(ob.trades))
	copy(result, ob.trades)
	return result
}
