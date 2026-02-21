package utils

import (
	"context"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/models"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/gofiber/contrib/websocket"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type BroadcastMessage struct {
	Type     string              `json:"type"`
	Tick     int                 `json:"tick"`
	Stocks   []models.StockPrice `json:"stocks,omitempty"`
	News     []models.News       `json:"news,omitempty"`
	IsActive bool                `json:"is_active"`
	IsPaused bool                `json:"is_paused"`
	Error    string              `json:"error,omitempty"`
}

type Client struct {
	Conn *websocket.Conn
	Send chan BroadcastMessage
}

type Hub struct {
	clients         map[*Client]bool
	broadcastSignal <-chan models.SimulationBroadcast
	store           *db.Store
	engine          *simulation.Engine
	Join            chan *Client
	Leave           chan *Client
}

func NewHub(tick <-chan models.SimulationBroadcast, store *db.Store, engine *simulation.Engine) *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		broadcastSignal: tick,
		store:           store,
		engine:          engine,
		Join:            make(chan *Client),
		Leave:           make(chan *Client),
	}
}

func (h *Hub) Run(ctx context.Context) {
	logs.Log.Info("WebSocket Hub running")
	for {
		select {
		case <-ctx.Done():
			logs.Log.Info("WebSocket Hub stopping via context")
			return

		case c := <-h.Join:
			h.clients[c] = true
			logs.Log.Debug("Client registered in hub", zap.Int("total_clients", len(h.clients)))

		case c := <-h.Leave:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
				logs.Log.Debug("Client deregistered from hub", zap.Int("total_clients", len(h.clients)))
			}

		case msg := <-h.broadcastSignal:
			tick := msg.Tick
			logs.Log.Debug("Hub received broadcast signal", zap.Int("tick", tick))

			// Hub now fetches the state data
			var prices []models.StockPrice
			var news []models.News
			var err error

			if msg.IsActive {
				prices, err = h.store.GetStocksAtTick(ctx, tick)
				if err != nil {
					logs.Log.Error("Hub failed to fetch prices for tick", zap.Int("tick", tick), zap.Error(err))
					if err == pgx.ErrNoRows {
						logs.Log.Warn("No prices found for tick, stopping simulation", zap.Int("tick", tick))
						h.engine.Stop()

						broadcastMsg := BroadcastMessage{
							Type:     "SIMULATION_ENDED",
							Tick:     tick,
							IsActive: false,
							IsPaused: false,
							Error:    "Simulation reached the end of available data",
						}
						h.broadcastToAll(broadcastMsg)
						continue // Skip normal update
					}
				}

				news, err = h.store.GetNewsAtTickForAll(ctx, tick)
				if err != nil {
					logs.Log.Error("Hub failed to fetch news for tick", zap.Int("tick", tick), zap.Error(err))
				}
			}

			msgType := "UPDATE"
			if !msg.IsActive {
				msgType = "SIMULATION_ENDED"
			}

			broadcastMsg := BroadcastMessage{
				Type:     msgType,
				Tick:     tick,
				Stocks:   prices,
				News:     news,
				IsActive: msg.IsActive,
				IsPaused: msg.IsPaused,
			}

			h.broadcastToAll(broadcastMsg)
		}
	}
}

func (h *Hub) broadcastToAll(msg BroadcastMessage) {
	count := 0
	for client := range h.clients {
		select {
		case client.Send <- msg:
			count++
		default:
			logs.Log.Debug("Skipped broadcast for slow client")
		}
	}
	// logs.Log.Info("Broadcasted message to clients", zap.String("type", msg.Type), zap.Int("tick", msg.Tick), zap.Int("clients_notified", count))
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		if err := c.Conn.WriteJSON(msg); err != nil {
			logs.Log.Debug("WebSocket write error, closing client connection", zap.Error(err))
			return
		}
	}
}

func (h *Hub) GetConnectionCount() int {
	return len(h.clients)
}
