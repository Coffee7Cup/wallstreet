package routers

import (
	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func AdminRouter(router fiber.Router, handler *controllers.AdminHandler) {
	logs.Log.Info("Initializing Admin Router")
	router.Post("/login", handler.AdminLoginHandler)

	engineGroup := router.Group("/engine", middleware.JWTMiddleware(), middleware.VerifyAdmin())
	engineGroup.Post("/start", handler.StartEngine)
	engineGroup.Post("/stop", handler.StopEngine)
	engineGroup.Post("/pause", handler.PauseEngine)
	engineGroup.Post("/resume", handler.ResumeEngine)

	adminDocs := router.Group("/monitor", middleware.JWTMiddleware(), middleware.VerifyAdmin())
	adminDocs.Get("/", monitor.New(monitor.Config{Title: "Wallstreet Engine Monitor"}))
	adminDocs.Get("/ws", handler.MonitorWS)

	statsGroup := router.Group("/stats", middleware.JWTMiddleware(), middleware.VerifyAdmin())
	statsGroup.Get("/", handler.GetStats)
	statsGroup.Get("/logs", handler.GetLogs)
	statsGroup.Get("/top-traders", handler.GetTopTraders)

	tradeGroup := router.Group("/trade", middleware.JWTMiddleware(), middleware.VerifyAdmin())
	tradeGroup.Get("/trades/:user_id", handler.GetTradesByUserID)
	tradeGroup.Get("/trades/:user_id/:symbol", handler.GetTradesBySymbol)
	tradeGroup.Get("/trades/:user_id/:symbol/:limit", handler.GetTradesBySymbolAndLimit)

	simulationGroup := router.Group("/simulation", middleware.JWTMiddleware(), middleware.VerifyAdmin())
	simulationGroup.Post("/reset", handler.ResetSimulation)
}
