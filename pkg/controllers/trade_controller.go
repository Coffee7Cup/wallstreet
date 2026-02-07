package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/coffee7cup/wallstreet/pkg/utils"
	"github.com/gofiber/contrib/websocket"
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

// WebSocketHandler handles real-time trade requests and market data broadcasts.
func (h *TradeHandler) WebSocketHandler(c *websocket.Conn) {
	client := &utils.Client{
		Conn: c,
		Send: make(chan utils.BroadcastMessage, 256),
	}

	h.hub.Join <- client
	logs.Log.Info("New client joined WebSocket hub", zap.String("remote_addr", c.RemoteAddr().String()))

	// Send initial state immediately
	state := h.engine.GetState()
	prices, _ := h.store.GetStocksAtTick(context.Background(), state.Tick)

	client.Send <- utils.BroadcastMessage{
		Type:     "INITIAL_STATE",
		Tick:     state.Tick,
		Stocks:   prices,
		IsActive: state.IsActive,
		IsPaused: state.IsPaused,
	}

	go client.WritePump()

	defer func() {
		h.hub.Leave <- client
		logs.Log.Info("Client left WebSocket hub", zap.String("remote_addr", c.RemoteAddr().String()))
	}()

	for {
		var tradeReq TradeRequest

		if err := c.ReadJSON(&tradeReq); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logs.Log.Warn("WebSocket read error", zap.Error(err))
			}
			return // client disconnected
		}

		// validate simulation state (time only)
		state := h.engine.GetState()
		if !state.IsActive || state.IsPaused {
			logs.Log.Debug("Trade request ignored: simulation not active or paused", zap.Any("tradeReq", tradeReq))
			continue
		}

		var trade models.Trade
		trade.UserID = tradeReq.UserID
		trade.CompanyID = tradeReq.CompanyID
		trade.TradeType = tradeReq.TradeType
		trade.Quantity = tradeReq.Quantity

		// Fetch price at current tick to get the simulation date for this company
		price, err := h.store.GetStockAtTickOfCompany(context.Background(), state.Tick, trade.CompanyID)
		if err == nil {
			trade.Date = price.Date
		} else {
			trade.Date = time.Now() // fallback
		}
		trade.Timestamp = time.Now()

		logs.Log.Info("Trade execution attempt", zap.Int("user_id", trade.UserID), zap.String("type", trade.TradeType), zap.Int("company_id", trade.CompanyID), zap.Int("quantity", trade.Quantity))

		// execute trade
		switch trade.TradeType {
		case "BUY":
			err = h.store.BuyStocks(context.Background(), trade)
		case "SELL":
			err = h.store.SellStocks(context.Background(), trade)
		default:
			logs.Log.Warn("Unknown trade type received", zap.String("type", trade.TradeType))
			err = fmt.Errorf("unknown trade type: %s", trade.TradeType)
		}

		if err != nil {
			logs.Log.Error("Trade execution failed", zap.Any("trade", trade), zap.Error(err))
			// send error back to client in a standardized way via the write pump (thread-safe)
			client.Send <- utils.BroadcastMessage{
				Type:  "ERROR",
				Error: err.Error(),
			}
			continue
		}

		logs.Log.Info("Trade execution successful", zap.Int("user_id", trade.UserID), zap.String("type", trade.TradeType))
	}
}
