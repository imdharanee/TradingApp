package models

import (
	"time"

	"github.com/google/uuid"
)

// OrderType represents the type of order (limit or market)
type OrderType string

// OrderSide represents the side of the order (buy or sell)
type OrderSide string

// OrderStatus represents the status of an order
type OrderStatus string

const (
	// Order types
	LimitOrder  OrderType = "LIMIT"
	MarketOrder OrderType = "MARKET"

	// Order sides
	BuyOrder  OrderSide = "BUY"
	SellOrder OrderSide = "SELL"

	// Order statuses
	Open      OrderStatus = "OPEN"
	Partial   OrderStatus = "PARTIAL"
	Filled    OrderStatus = "FILLED"
	Cancelled OrderStatus = "CANCELLED"
)

// Order represents an order in the exchange
type Order struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Type      OrderType   `json:"type"`
	Side      OrderSide   `json:"side"`
	Status    OrderStatus `json:"status"`
	Price     float64     `json:"price"`
	Quantity  float64     `json:"quantity"`
	Filled    float64     `json:"filled"`
	Remaining float64     `json:"remaining"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// NewOrder creates a new order with a unique ID
func NewOrder(userID string, orderType OrderType, side OrderSide, price, quantity float64) *Order {
	now := time.Now()
	return &Order{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      orderType,
		Side:      side,
		Status:    Open,
		Price:     price,
		Quantity:  quantity,
		Filled:    0,
		Remaining: quantity,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Fill updates the order when it is filled (partially or completely)
func (o *Order) Fill(quantity float64) {
	o.Filled += quantity
	o.Remaining -= quantity
	o.UpdatedAt = time.Now()

	if o.Remaining <= 0 {
		o.Status = Filled
	} else {
		o.Status = Partial
	}
}

// Cancel marks the order as cancelled
func (o *Order) Cancel() {
	o.Status = Cancelled
	o.UpdatedAt = time.Now()
}
