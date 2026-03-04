package controllers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/coffee7cup/wallstreet/pkg/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TradeHandler struct {
	store  *db.Store
	engine *simulation.Engine
	hub    *utils.Hub
}

func NewTradeHandler(
	store *db.Store,
	engine *simulation.Engine,
	hub *utils.Hub,
) *TradeHandler {
	return &TradeHandler{
		store:  store,
		engine: engine,
		hub:    hub,
	}
}

type TradeRequest struct {
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	TradeType string `json:"trade_type"`
	Quantity  int    `json:"quantity"`
}

// GetUserPortfolio fetches the stock portfolio of the authenticated user.
func (h *TradeHandler) GetUserPortfolio(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("Portfolio fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	state := h.engine.GetState()
	portfolio, err := h.store.GetUserPortfolio(c.Context(), id, state.Tick)
	if err != nil {
		logs.Log.Error("Portfolio fetch failed: internal error", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
			"err":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "user portfolio fetched successfully",
		"portfolio": portfolio,
	})
}

// GetUserTrades fetches the trade history of the authenticated user.
func (h *TradeHandler) GetUserTrades(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("Trades fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	trades, err := h.store.GetUserTrades(c.Context(), id)
	if err != nil {
		logs.Log.Error("Trades fetch failed: internal error", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user trades fetched successfully",
		"trades":  trades,
	})
}

// WebSocketHandler handles real-time trade requests and market data broadcasts.
func (h *TradeHandler) WebSocketHandler(c *websocket.Conn) {
	userID, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("WebSocket connection rejected: unauthorized")
		c.WriteJSON(utils.BroadcastMessage{
			Type:  "ERROR",
			Error: "unauthorized",
		})
		c.Close()
		return
	}

	client := &utils.Client{
		Conn: c,
		Send: make(chan utils.BroadcastMessage, 256),
	}

	h.hub.Join <- client
	logs.Log.Info("New client joined WebSocket hub", zap.Int("user_id", userID), zap.String("remote_addr", c.RemoteAddr().String()))

	go client.WritePump()

	defer func() {
		h.hub.Leave <- client
		logs.Log.Info("Client left WebSocket hub", zap.Int("user_id", userID), zap.String("remote_addr", c.RemoteAddr().String()))
	}()

	// Send initial state immediately
	state := h.engine.GetState()
	prices := h.store.CacheCurrentStocks
	if len(prices) == 0 {
		prices, _ = h.store.GetStocksAtTick(context.Background(), state.Tick)
	}

	client.Send <- utils.BroadcastMessage{
		Type:     "INITIAL_STATE",
		Tick:     state.Tick,
		Stocks:   prices,
		IsActive: state.IsActive,
		IsPaused: state.IsPaused,
	}

	if !state.IsActive {
		client.Send <- utils.BroadcastMessage{
			Type:  "SIMULATION_ENDED",
			Error: "simulation complete",
		}
		// Give time for message delivery
		time.Sleep(200 * time.Millisecond)
		return
	}

	for {
		var tradeReq TradeRequest

		if err := c.ReadJSON(&tradeReq); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logs.Log.Warn("WebSocket read error", zap.Error(err))
			}
			return // client disconnected
		}

		// Basic input validation
		if tradeReq.Quantity <= 0 {
			client.Send <- utils.BroadcastMessage{Type: "ERROR", Error: "quantity must be greater than zero"}
			continue
		}
		if tradeReq.TradeType != "BUY" && tradeReq.TradeType != "SELL" {
			client.Send <- utils.BroadcastMessage{Type: "ERROR", Error: "invalid trade type"}
			continue
		}

		var err error

		// Validate simulation state
		state := h.engine.GetState()
		if !state.IsActive || state.IsPaused {
			logs.Log.Debug("Trade request ignored: simulation not active or paused", zap.Any("tradeReq", tradeReq))
			continue
		}

		var trade models.Trade
		trade.UserID = userID // SECURE: Use verified userID from token
		trade.CompanyID = tradeReq.CompanyID
		trade.TradeType = tradeReq.TradeType
		trade.Quantity = tradeReq.Quantity

		// Fetch price at current tick to get the simulation date for this company
		// Optimized: Check cache first
		var price models.StockPrice
		var priceFound bool

		for _, p := range h.store.CacheCurrentStocks {
			if p.CompanyID == trade.CompanyID {
				price = p
				priceFound = true
				break
			}
		}

		if !priceFound {
			var err error
			price, err = h.store.GetStockAtTickOfCompany(context.Background(), state.Tick, trade.CompanyID)
			if err != nil {
				logs.Log.Error("Failed to get stock at tick for trade", zap.Error(err))
				client.Send <- utils.BroadcastMessage{
					Type:  "ERROR",
					Error: "stock price not found for current tick",
				}
				continue
			}
		}

		trade.Date = price.Date
		trade.Timestamp = time.Now()

		logs.Log.Info("Trade execution attempt", zap.Int("user_id", trade.UserID), zap.String("type", trade.TradeType), zap.Int("company_id", trade.CompanyID))

		// execute trade
		switch trade.TradeType {
		case "BUY":
			err = h.store.BuyStocks(context.Background(), trade)
		case "SELL":
			err = h.store.SellStocks(context.Background(), trade)
		default:
			err = fmt.Errorf("unknown trade type: %s", trade.TradeType)
		}

		if err != nil {
			logs.Log.Error("Trade execution failed", zap.Any("trade", trade), zap.Error(err))
			client.Send <- utils.BroadcastMessage{
				Type:  "ERROR",
				Error: err.Error(),
			}
			continue
		}

		logs.Log.Info("Trade execution successful", zap.Int("user_id", trade.UserID), zap.String("type", trade.TradeType))

		// Fetch fresh data for the user to update the UI immediately
		updatedPortfolio, pErr := h.store.GetUserPortfolio(context.Background(), userID, state.Tick)
		updatedUser, uErr := h.store.GetUserByID(context.Background(), userID)

		if pErr == nil && uErr == nil {
			client.Send <- utils.BroadcastMessage{
				Type:      "TRADE_UPDATE",
				Tick:      state.Tick,
				Portfolio: updatedPortfolio,
				Balance:   updatedUser.CashBalance,
			}
		} else {
			logs.Log.Error("Failed to fetch fresh user data after trade", zap.Error(pErr), zap.Error(uErr))
		}
	}
}

func (h *TradeHandler) GetUserTradesLimit(c *fiber.Ctx) error {
	id, ok := c.Locals("user_id").(int)
	if !ok {
		logs.Log.Warn("Trades fetch failed: unauthorized")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	search := c.Query("search")

	trades, err := h.store.GetUserTradesSearchLimit(c.Context(), id, limit, offset, search)
	if err != nil {
		logs.Log.Error("Trades fetch failed: internal error", zap.Int("user_id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user trades fetched successfully",
		"trades":  trades,
	})
}
