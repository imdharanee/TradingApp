package api

import (
	"net/http"

	"dex/internal/engine"
	"dex/internal/models"
	"dex/internal/orderbook"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the DEX API
type Handler struct {
	orderBook      *orderbook.OrderBook
	matchingEngine *engine.MatchingEngine
}

// NewHandler creates a new API handler
func NewHandler(orderBook *orderbook.OrderBook, matchingEngine *engine.MatchingEngine) *Handler {
	return &Handler{
		orderBook:      orderBook,
		matchingEngine: matchingEngine,
	}
}

// SetupRoutes sets up the Gin routes
func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.POST("/orders", h.CreateOrder)
	router.GET("/orders", h.ListOrders)
	router.GET("/orders/:id", h.GetOrder)
	router.GET("/orderbook", h.GetOrderBook)
	router.GET("/trades", h.GetTrades)
}

// OrderRequest represents a request to create a new order
type OrderRequest struct {
	UserID   string           `json:"user_id"`
	Type     models.OrderType `json:"type"`
	Side     models.OrderSide `json:"side"`
	Price    float64          `json:"price"`
	Quantity float64          `json:"quantity"`
}

// CreateOrder handles order creation
func (h *Handler) CreateOrder(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate request
	if req.UserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameters: user_id is required"})
		return
	}

	if req.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameters: quantity must be greater than 0"})
		return
	}

	if req.Type == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameters: type is required (LIMIT or MARKET)"})
		return
	}

	if req.Side == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameters: side is required (BUY or SELL)"})
		return
	}

	if req.Type == models.LimitOrder && req.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order parameters: price must be greater than 0 for LIMIT orders"})
		return
	}

	// Create the order
	order := models.NewOrder(req.UserID, req.Type, req.Side, req.Price, req.Quantity)

	// Submit to matching engine
	h.matchingEngine.SubmitOrder(order)

	// Return the created order
	c.JSON(http.StatusOK, order)
}

// ListOrders handles listing all orders
func (h *Handler) ListOrders(c *gin.Context) {
	// List all orders (simplified - in a real system you'd want pagination)
	buys := h.orderBook.GetBuyOrders()
	sells := h.orderBook.GetSellOrders()

	// Combine buy and sell orders
	allOrders := append(buys, sells...)

	c.JSON(http.StatusOK, allOrders)
}

// GetOrder handles retrieving a specific order by ID
func (h *Handler) GetOrder(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, exists := h.orderBook.GetOrder(orderID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetOrderBook handles retrieving the current state of the order book
func (h *Handler) GetOrderBook(c *gin.Context) {
	// Get current order book state
	response := gin.H{
		"buys":  h.orderBook.GetBuyOrders(),
		"sells": h.orderBook.GetSellOrders(),
	}

	c.JSON(http.StatusOK, response)
}

// GetTrades handles retrieving recent trades
func (h *Handler) GetTrades(c *gin.Context) {
	// Get recent trades
	trades := h.orderBook.GetTrades()

	c.JSON(http.StatusOK, trades)
}
