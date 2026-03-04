package routers

import (
	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func VerifyWebSocketUpgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}

func TradeRouter(router fiber.Router, tradeHandler *controllers.TradeHandler) {
	logs.Log.Info("Initializing Trade Router")

	// Portfolio & Private Trade Routes
	router.Get("/portfolio", middleware.JWTMiddleware(), tradeHandler.GetUserPortfolio)
	router.Get("/trades", middleware.JWTMiddleware(), tradeHandler.GetUserTradesLimit)

	// WebSocket Route
	router.Get("/ws", middleware.JWTMiddleware(), VerifyWebSocketUpgrade, websocket.New(tradeHandler.WebSocketHandler))
}
