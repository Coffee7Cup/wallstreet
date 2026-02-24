package routers

import (
	"time"

	"github.com/coffee7cup/wallstreet/pkg/controllers"
	"github.com/coffee7cup/wallstreet/pkg/logs"
	"github.com/coffee7cup/wallstreet/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func MarketRouter(router fiber.Router, handler *controllers.MarketHandler) {
	logs.Log.Info("Initializing Market Router")

	// Cache middleware for static data (5 minutes)
	cacheMiddleware := cache.New(cache.Config{
		Expiration:   5 * time.Minute,
		CacheControl: true,
	})

	router.Get("/companies", cacheMiddleware, handler.GetAllCompanies)
	router.Get("/sectors", cacheMiddleware, handler.GetAllSectors)
	router.Get("/companies/:id/peers", cacheMiddleware, handler.GetPeerComparison)
	router.Get("/companies/:id", cacheMiddleware, handler.GetCompany)
	router.Get("/profit-loss/:id", cacheMiddleware, handler.GetProfitLoss)
	router.Get("/balance-sheets/:id", cacheMiddleware, handler.GetBalanceSheets)
	router.Get("/cash-flows/:id", cacheMiddleware, handler.GetCashFlows)
	router.Get("/ratios/:id", cacheMiddleware, handler.GetRatios)

	// History and dynamic data shouldn't be aggressively cached to preserve websocket timing
	router.Get("/stocks-history/:id", handler.GetStocksTillTick)
	router.Get("/stocks-history", handler.GetStocksTillTick)

	router.Get("/news", middleware.JWTMiddleware(), handler.GetNewsTillTickLimitSearch)
	router.Get("/news/sector", middleware.JWTMiddleware(), handler.GetNewsBySectorTillTick)

	router.Get("/drawings/:id", handler.GetDrawings)
	router.Post("/drawings/:id", handler.SaveDrawings)
}
