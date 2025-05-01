package models

import (
	"time"

	"github.com/google/uuid"
)

// Trade represents a completed trade between two orders
type Trade struct {
	ID          string    `json:"id"`
	BuyOrderID  string    `json:"buy_order_id"`
	SellOrderID string    `json:"sell_order_id"`
	Price       float64   `json:"price"`
	Quantity    float64   `json:"quantity"`
	ExecutedAt  time.Time `json:"executed_at"`
}

// NewTrade creates a new trade with a unique ID
func NewTrade(buyOrderID, sellOrderID string, price, quantity float64) *Trade {
	return &Trade{
		ID:          uuid.New().String(),
		BuyOrderID:  buyOrderID,
		SellOrderID: sellOrderID,
		Price:       price,
		Quantity:    quantity,
		ExecutedAt:  time.Now(),
	}
}
