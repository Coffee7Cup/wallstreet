package controllers

import (
	"fmt"
	"os"
	"time"

	"github.com/coffee7cup/wallstreet/pkg/db"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/coffee7cup/wallstreet/pkg/simulation"
	"github.com/coffee7cup/wallstreet/pkg/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"go.uber.org/zap"
)

type AdminHandler struct {
	store  *db.Store
	engine *simulation.Engine
	hub    *utils.Hub
}

func NewAdminHandler(store *db.Store, engine *simulation.Engine, hub *utils.Hub) *AdminHandler {
	return &AdminHandler{
		store:  store,
		engine: engine,
		hub:    hub,
	}
}

// AdminLoginHandler handles admin authentication.
func (h *AdminHandler) AdminLoginHandler(c *fiber.Ctx) error {
	type Body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var body Body
	if err := c.BodyParser(&body); err != nil {
		logs.Log.Warn("Admin login failed: invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	logs.Log.Info("Admin login attempt", zap.String("username", body.Username), zap.String("email", body.Email))

	admin, err := h.store.GetAdmin(c.Context(), body.Username, body.Email)
	if err != nil {
		logs.Log.Warn("Admin login failed: invalid credentials", zap.String("username", body.Username), zap.Error(err))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := middleware.GenerateToken(admin.ID, admin.Username, true)
	if err != nil {
		logs.Log.Error("Admin login failed: token generation error", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	logs.Log.Info("Admin login successful", zap.String("username", admin.Username))

	return c.JSON(fiber.Map{
		"message": "admin login successful",
		"admin": fiber.Map{
			"id":       admin.ID,
			"username": admin.Username,
			"email":    admin.Email,
		},
		"token": token,
	})
}

// StartEngine activates the simulation engine.
func (h *AdminHandler) StartEngine(c *fiber.Ctx) error {
	type Request struct {
		StartTick int `json:"start_tick"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		logs.Log.Warn("Start engine failed: invalid format", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid format"})
	}

	logs.Log.Warn("Starting simulation engine", zap.Int("start_tick", req.StartTick), zap.String("category", logs.CategoryEngine))
	h.engine.Start(req.StartTick)

	return c.JSON(fiber.Map{"message": "engine started"})
}

// StopEngine stops the simulation engine.
func (h *AdminHandler) StopEngine(c *fiber.Ctx) error {
	logs.Log.Warn("Stopping simulation engine", zap.String("category", logs.CategoryEngine))
	h.engine.Stop()
	return c.JSON(fiber.Map{"message": "engine stopped"})
}

// PauseEngine pauses the simulation engine.
func (h *AdminHandler) PauseEngine(c *fiber.Ctx) error {
	logs.Log.Warn("Pausing simulation engine", zap.String("category", logs.CategoryEngine))
	h.engine.Pause()
	return c.JSON(fiber.Map{"message": "engine paused"})
}

// ResumeEngine resumes the simulation engine.
func (h *AdminHandler) ResumeEngine(c *fiber.Ctx) error {
	logs.Log.Warn("Resuming simulation engine", zap.String("category", logs.CategoryEngine))
	h.engine.Resume()
	return c.JSON(fiber.Map{"message": "engine resumed"})
}

// GetStats returns system statistics.
func (h *AdminHandler) GetStats(c *fiber.Ctx) error {
	stats := h.engine.GetState()
	return c.JSON(fiber.Map{
		"active_connections": h.hub.GetConnectionCount(),
		"active_requests":    middleware.GetActiveRequests(),
		"simulation_state":   stats,
	})
}

// GetLogs returns the last 100 lines of the engine log.
func (h *AdminHandler) GetLogs(c *fiber.Ctx) error {
	file, err := os.Open("./logs/engine.log")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not open log file"})
	}
	defer file.Close()

	// Simple tail implementation (read last 10KB)
	stat, err := file.Stat()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not open log file"})
	}
	size := stat.Size()
	limit := int64(10240) // 10KB
	if size < limit {
		limit = size
	}

	buffer := make([]byte, limit)
	_, err = file.ReadAt(buffer, size-limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not read log file"})
	}

	return c.SendString(string(buffer))
}

// MonitorWS streams system stats over WebSocket.
func (h *AdminHandler) MonitorWS(c *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return fiber.ErrUpgradeRequired
	}

	return websocket.New(func(conn *websocket.Conn) {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			cUsage, _ := cpu.Percent(0, false)
			mUsage, _ := mem.VirtualMemory()

			stats := fiber.Map{
				"cpu_usage":          fmt.Sprintf("%.2f%%", cUsage[0]),
				"ram_usage":          fmt.Sprintf("%.2f GB / %.2f GB", float64(mUsage.Used)/1e9, float64(mUsage.Total)/1e9),
				"active_connections": h.hub.GetConnectionCount(),
				"active_requests":    middleware.GetActiveRequests(),
				"simulation_tick":    h.engine.GetState().Tick,
				"is_active":          h.engine.GetState(),
				"ts":                 time.Now().Format(time.RFC3339),
			}

			if err := conn.WriteJSON(stats); err != nil {
				return
			}
		}
	})(c)
}
